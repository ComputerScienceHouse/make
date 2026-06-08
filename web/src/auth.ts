import { reactive } from "vue";

export interface UserInfo {
    uuid: string;
    email: string;
    preferred_username: string;
    name: string;
    groups: string[];
}

export const authState = reactive({
    user: null as UserInfo | null,
});

export async function loadUser() {

    const res = await fetch("/api/me", {
        credentials: "include",
    });

    if (!res.ok) {
        authState.user = null;
        return;
    }

    authState.user = await res.json();

}

export function logout() {
    authState.user = null;
}