<script context="module" lang="ts">
    import {getPackage, getPackageTags, getUser, Package_Base, Tag_Base, User_Base} from "$lib/api";

    export async function load({ params }) {
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
            stuff: {
                pack: pack,
                author: author,
                tags: tags,
            }
        }
    }
</script>

<slot />