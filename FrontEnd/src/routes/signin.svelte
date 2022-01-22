<script context="module" lang="ts">
    export async function handle({ event, resolve }) {
        resolve.ssr = false;
        return await resolve(event);
    }
</script>

<script lang="ts">
    import Paper, {Content} from "@smui/paper"
    import {Icon} from "@smui/common"
    import {page} from "$app/stores";
    import {onMount} from "svelte";
    import {googleClientId} from "$lib/auth";

    let signInFailed = $page.url.searchParams.get('failure')
    let googleSignInBtn : HTMLElement
    let form : HTMLFormElement

    onMount(() => {
        // eslint-disable-next-line no-undef
        google.accounts.id.renderButton(googleSignInBtn, {
            theme: "outline",
            text: "signin_with",
            size: "large",
            logo_alignment: "left",
            shape: "rectangular",
        })
    })
</script>

<div id="g_id_onload"
     data-client_id="{googleClientId}"
     data-login_uri="http://localhost:8000/oauth?redirect={$page.url.searchParams.get('redirect') || encodeURI('http://localhost:3000')}&failure-redirect=${encodeURI('http://localhost:3000/signin')}"
     data-auto_prompt="false">
</div>

<h1>Sign-In/Sign-Up</h1>
<form bind:this={form}>
    {#if signInFailed}
        <Paper color="primary">
            <Content class="flex items-center item">
                <Icon class="material-icons mr-5">warning</Icon>
                <p>
                    Failed to sign in: {signInFailed}
                </p>
            </Content>
        </Paper>
    {/if}

    <div bind:this={googleSignInBtn}></div>
</form>

<style>
    .info div > * {
        display: inline-block;
    }
</style>