<script lang="ts">
    import Textfield from "@smui/textfield";
    import HelperText from "@smui/textfield/helper-text";
    import Button, {Label} from '@smui/button';
    import PackageEdit from "$lib/components/PackageEdit.svelte";
    import {Package_Base, Tag_Base} from "$lib/api.ts";

    let packageStruct : Package_Base = new Package_Base()
    let packageTags : Tag_Base[] = []

    async function createPackage() {
        console.log("Create package!")
    }
</script>

<h1>Create new Mod</h1>
<Textfield bind:value={packageStruct.name} label="Name">
    <HelperText slot="helper">Name of you package that can be used as reference. Only a-z, A-Z, 0-9, - and _</HelperText>
</Textfield>
<PackageEdit pack={packageStruct} tags={packageTags} />
<Button variant="outlined" on:click={createPackage}>
    <Label>Create Package</Label>
</Button>
