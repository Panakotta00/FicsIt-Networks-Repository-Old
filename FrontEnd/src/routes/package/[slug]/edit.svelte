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
    import Button, {Label} from '@smui/button';
    import PackageEdit from "$lib/components/PackageEdit.svelte";

    async function updatePackage() {
        console.log("Update package!")
    }
</script>

<PackageEdit pack={packageStruct} tags={packageTags} />
<Button variant="outlined" on:click={updatePackage}>
    <Label>Update Package</Label>
</Button>
