<script lang="ts">
	import "../app.css";
    import {page} from "$app/stores";
    import {getCookie} from "$lib/util";
    import {onMount} from "svelte";
    import {googleClientId} from "$lib/auth";
    import {cacheExchange, dedupExchange, fetchExchange, initClient, Operation} from '@urql/svelte';
    import { authExchange, AuthConfig } from '@urql/exchange-auth';

    let hasLoginChecked = false
    let isLoggedIn = false

    initClient({
        url: 'http://localhost:8000/query',
        exchanges: [
            dedupExchange,
            cacheExchange,
            authExchange(<AuthConfig<string>>{
                addAuthToOperation(params: { authState: string|null; operation: Operation }): Operation {
                    params.operation.context.fetchOptions = {
                        headers: {
                            Authorization: params.authState
                        }
                    }
                    return params.operation
                },
                getAuth: async () : Promise<string|null> => {
                    return (typeof document !== "undefined" && getCookie("token")) || ""
                },
                didAuthError({ error }): boolean {
                    return error.message.indexOf('user not logged in') >= 0;
                }
            }),
            fetchExchange,
        ]
    });

    onMount(() => {
        isLoggedIn = getCookie("token") != null
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
