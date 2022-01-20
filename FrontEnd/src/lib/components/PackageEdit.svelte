<script lang="ts">
    import Textfield from "@smui/textfield";
    import HelperText from "@smui/textfield/helper-text";
    import Card from "@smui/card";
    import TagList from "$lib/components/TagList.svelte";
    import {Package_Base, Tag_Base} from "$lib/api";

    export let pack : Package_Base = null
    export let tags : Tag_Base[] = null
</script>

<Textfield bind:value={pack.displayName} label="Display-Name">
    <HelperText slot="helper">Human-Readable Display-Name of your package that is shown in UI</HelperText>
</Textfield>
<div class="margins">
    <Textfield
            style="width: 100%;"
            helperLine$style="width: 100%;"
            textarea
            bind:value={pack.description}
            label="Description"
    >
        <HelperText slot="helper">Markdown formatted description of your package</HelperText>
    </Textfield>
</div>
<Textfield bind:value={pack.sourceLink} label="Source-URL">
    <HelperText slot="helper">Optional URL to your source repository like GitHub-Repository</HelperText>
</Textfield>
<Card padded class="m-3 w-100">
    <TagList bind:tags={tags} editable=true />
</Card>
