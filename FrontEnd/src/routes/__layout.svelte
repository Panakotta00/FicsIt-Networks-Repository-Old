<script>
	import "../app.css";
    import {page} from "$app/stores";
    import {getCookie} from "$lib/util";
    import {onMount} from "svelte";
    import {googleClientId} from "$lib/auth";

    let hasLoginChecked = false
    let isLoggedIn = false

    onMount(() => {
        isLoggedIn = getCookie("g_state") != null
        hasLoginChecked = true
        console.log("login: ", isLoggedIn)

        if (!isLoggedIn && $page.url.pathname !== "/signin") {
            // eslint-disable-next-line no-undef
            google.accounts.id.prompt();
        }
    })
</script>

{#if $page.url.pathname !== "/signin"}
    <div id="g_id_onload"
         data-client_id="{googleClientId}"
         data-login_uri="http://localhost:8000/oauth?redirect={encodeURI($page.url.href)}&failure-redirect={encodeURI('http://localhost:3000/signin')}"
         data-auto_prompt="false">
    </div>
{/if}

<slot />
