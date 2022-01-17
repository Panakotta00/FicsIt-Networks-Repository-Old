<script context="module" lang="ts">
    import {getPackage, getPackageTags, getUser, Package_Base, Tag_Base, User_Base} from "$lib/api";

    export async function load({ params, fetch }) {
        var pack : Package_Base|null = null
        var author : User_Base|null = null
        var tags : Tag_Base[]|null = null

        try {
            pack = await getPackage(params["slug"])
            let authorP = getUser(pack.creatorId);
            let tagsP = getPackageTags(pack.id);

            author = await authorP
            tags = await tagsP
        } catch {
            return {
                status: 404,
                error: `Package with Ref '${params["slug"]}' not found`
            }
        }

        if (params["slug"].match(/\d+/)) {
            return {
                status: 307,
                redirect: `${pack.name}`
            }
        }

        return {
            props: {
                packageStruct: pack,
                authorStruct: author,
                packageTags: tags,
            }
        }
    }
</script>

<script lang="ts">
    import TagList from "../../components/TagList.svelte";

    export let packageStruct : Package_Base = null
    export let authorStruct : User_Base = null
    export let packageTags : Tag_Base[] = null
</script>

<h1>{packageStruct.displayName}</h1>
<p>{packageStruct.description}</p>
{#if packageStruct.sourceLink}
    <p>Source: <a href="{packageStruct.sourceLink}">{packageStruct.sourceLink}</a></p>
{/if}
<p>By: <a href="/user/{authorStruct.name}">{authorStruct.name}</a></p>
<TagList bind:tags={packageTags} editable=true />
<TagList bind:tags={packageTags} />
