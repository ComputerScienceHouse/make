<template>
    <nav class="navbar fixed-top navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            
            <RouterLink class="navbar-brand" to="/">
                <img class="object-fit-contain ms-2 me-2" height="32" src="https://assets.csh.rit.edu/pubsite/csh_logo_square.svg">
                Make
            </RouterLink>
            <button class="navbar-toggler navbar-toggler-right" type="button" data-bs-toggle="collapse"
                data-bs-target="#navbarResponsive">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarResponsive">
                <ul class="navbar-nav mr-auto">
                    <RouterLink class="nav-link" to="/history">History</RouterLink>
                </ul>
                <UserMenu />
            </div>
        </div>
    </nav>
</template>