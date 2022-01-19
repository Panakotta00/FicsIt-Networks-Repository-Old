<script context="module" lang="ts">
    import {Package_Base, Tag_Base, User_Base} from "$lib/api";

    let packageStruct : Package_Base = null
    let authorStruct : User_Base = null
    let packageTags : Tag_Base[] = null

    export async function load({ params, fetch, session, stuff }) {
        packageStruct = stuff.pack
        authorStruct = stuff.author
        packageTags = stuff.tags
        return true
    }
</script>
<script>
    import TagList from "$lib/components/TagList.svelte";
</script>

<h1>{packageStruct.displayName}</h1>
<p>{packageStruct.description}</p>
{#if packageStruct.sourceLink}
    <p>Source: <a href="{packageStruct.sourceLink}">{packageStruct.sourceLink}</a></p>
{/if}
<p>By: <a href="/user/{authorStruct.name}">{authorStruct.name}</a></p>
<TagList bind:tags={packageTags} />
