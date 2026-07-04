package worker

import (
	"errors"
	"log"
	"makedotcsh/database"
	"makedotcsh/ldap"
	"makedotcsh/models"
	"slices"
	"time"

	goldap "github.com/go-ldap/ldap/v3"
)

var WorkerTrigger = make(chan struct{})

func TriggerWorker() {
	WorkerTrigger <- struct{}{}
}

func StartWorker(interval time.Duration, manual <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// initial run
		worker()
		log.Printf("[WORKER] Continuing in %s", interval.String())

		for {
			select {
			case <-ticker.C:
				worker()
				log.Printf("[WORKER] Continuing in %s", interval.String())
			case <-manual:
				worker()
				log.Printf("[WORKER] manual trigger")
			}
		}

	}()
}

func workerErr(err error) {
	if err != nil {
		log.Panicf("[WORKER] Error while handling job: \n %v", err)
	}
}

func worker() {
	log.Println("[WORKER] Starting job...")

	areas, err := database.Helper.GetAllAreas()
	workerErr(err)

	var ldapGroups []string

	for _, area := range areas {
		if area.LdapGroup == "" {
			continue
		}

		ldapGroups = append(ldapGroups, area.LdapGroup)
	}

	areaMap := map[int]models.Area{}
	for _, area := range areas {
		areaMap[area.ID] = area
	}

	groupAreaMap := map[string]int{}
	for _, area := range areas {
		if area.LdapGroup != "" {
			groupAreaMap[area.LdapGroup] = area.ID
		}
	}

	// pass 1, ensure everyone with an existing ldap role should
	// actually have it (ex training expired or was revoked)
	log.Println("[WORKER] Checking validity of all ldap group holders")
	for _, group := range ldapGroups {
		members, err := ldap.GetDirectGroupMembers(group)
		workerErr(err)

		//log.Printf("[WORKER] group: %s members: %s", group, members)

		for _, member := range members {
			userAreas, err := database.Helper.GetAllAreasWithUserAccess(member.UUID)
			workerErr(err)

			if !slices.Contains(userAreas, groupAreaMap[group]) {
				log.Printf("[WORKER] Removing user %s from group %s\n",
					member.Username, group)

				err = ldap.RemoveUserFromGroup(member.Username, group)
				workerErr(err)
			}
		}
	}
	log.Println("[WORKER] Completed validity check")

	// pass 2, ensure everyone with a user training has the
	// the ldap role for all areas where they have met all
	// criteria

	log.Println("[WORKER] Issuing new ldap groups")
	rows, err := database.DB.Query("SELECT DISTINCT user_uuid FROM user_trainings")
	workerErr(err)

	for rows.Next() {
		var uuid string

		err = rows.Scan(&uuid)
		workerErr(err)

		areaIds, err := database.Helper.GetAllAreasWithUserAccess(uuid)
		workerErr(err)

		// no need to worry about removing ldap roles, should've been
		// done in the first step if necessary
		// add ldap role

		username, err := ldap.UUIDtoUsername(uuid)
		workerErr(err)

		for _, id := range areaIds {
			area, ok := areaMap[id]
			if !ok {
				workerErr(errors.New("area does not exist in map"))
			}

			if area.LdapGroup == "" {
				continue
			}

			log.Printf("[WORKER] Adding user %s to group %s (area: %d)\n",
				username, area.LdapGroup, area.ID)
			err = ldap.AddUserToGroup(username, area.LdapGroup)

			// ignore error is it is "atribute already exists"
			if ldapErr, ok := err.(*goldap.Error); ok &&
				ldapErr.ResultCode == goldap.LDAPResultAttributeOrValueExists {
				continue
			}

			workerErr(err)
		}
	}

	log.Printf("[WORKER] Done\n")
}
