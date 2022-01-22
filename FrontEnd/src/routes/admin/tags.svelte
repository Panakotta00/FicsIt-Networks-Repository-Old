<script lang="ts">
    import Accordion, {Panel, Header, Content} from '@smui-extra/accordion';
    import {onMount} from "svelte";
    import {getTags, Tag, Tag_Base} from "$lib/api";
    import Textfield from "@smui/textfield";
    import HelperText from "@smui/textfield/helper-text";
    import Button, {Label, Icon} from "@smui/button";
    import Snackbar, {SnackbarComponentDev} from "@smui/snackbar";
    import IconButton, { Icon as ButtonIcon } from '@smui/icon-button';

    let tags : Tag_Base[] = []
    let accordion : Accordion
    let panels = {}
    let nameFields = {}
    let snackbarTagChangeSavedText = ""
    let snackbarTagChangeSaved : SnackbarComponentDev
    let tagNegativeID = -1

    function newTag() {
        if (!tags.find(tag => tag.name == "New Tag")) {
            const tag = new Tag({id: tagNegativeID-- as bigint, name: "New Tag"})
            tags.push(tag)
            tags = tags
            setTimeout(() => {
                panels[tag.id].setOpen(true)
                var field = nameFields[tag.id]
                field.focus()
                let input = field.getElement().querySelectorAll("input")[0] as HTMLInputElement
                input.select()
            }, 0)
        } else {
            panels[tags[tags.length-1].id].setOpen(true)
            nameFields[tags[tags.length-1].id].focus()
        }
    }

    function tagChange(tag : Tag_Base) {
        // ignore "New Tag"
        if (tag.name == "New Tag") return;

        if (tag.id < 0) {
            // Create new tag & update tag.id with new DB id or re-fetch all tags
        } else {
            // Update existing tag
        }

        snackbarTagChangeSavedText = `Tag '${tag.name}' saved!`
        snackbarTagChangeSaved.open()
        setTimeout(() => snackbarTagChangeSaved.close(), 2000)
    }

    function deleteTag(tag) {
        if (tag.name == "New Tag") {
            // Remove tag

        }

        // Remove tag animation
        let panelRemoveAnimation = () => {
            let panel = panels[tag.id].getElement()
            let startHeight = panel.scrollHeight
            panel.classList.add("smui-accordion__panel--removed")
            panel.style.height = startHeight + 'px';
            requestAnimationFrame(function() {
                panel.style.height = 0 + 'px';
            });
            panel.addEventListener("transitionend", e => {
                if (e.propertyName == "height") {
                    panel.classList.remove("smui-accordion__panel--removed")
                    panel.style.height = "auto"
                    tags = tags.filter(t => {
                        return t.id != tag.id
                    })
                }
            })
        }

        console.log(panels[123])
        let isPanelOpen = false
        for (let key in panels) {
            let panelP = panels[key]
            if (panelP?.isOpen()) {
                panelP.setOpen(false)
                if (!isPanelOpen) {
                    panelP.getElement().addEventListener("SMUIAccordionPanel:closed", () => {
                        panelRemoveAnimation()
                    }, {once: true})
                }
                isPanelOpen = true
            }
        }
        if (!isPanelOpen) {
            panelRemoveAnimation()
        }

        snackbarTagChangeSavedText = `Tag '${tag.name}' removed!`
        snackbarTagChangeSaved.open()
        setTimeout(() => snackbarTagChangeSaved.close(), 2000)
    }

    onMount(async () => {
        tags = await getTags() as Tag_Base[]
    })
</script>

<Accordion bind:this={accordion}>
    {#each tags as tag}
        <Panel bind:this={panels[tag.id]}>
            <Header>
                {tag.name}
                <IconButton slot="icon" on:click={e => {e.stopPropagation(); deleteTag(tag)}}>
                    <ButtonIcon class="material-icons">delete_forever</ButtonIcon>
                </IconButton>
            </Header>
            <Content>
                <Textfield bind:value={tag.name} label="Tag-Name" bind:this={nameFields[tag.id]} on:change={() => tagChange(tag)}>
                    <HelperText slot="helper">Human-Readable name of the tag that is shown in UI</HelperText>
                </Textfield>
                <Textfield
                        style="width: 100%;"
                        helperLine$style="width: 100%;"
                        textarea
                        bind:value={tag.description}
                        label="Tag-Description"
                        on:change={() => tagChange(tag)}
                >
                    <HelperText slot="helper">Markdown formatted description of the tag</HelperText>
                </Textfield>
            </Content>
        </Panel>
    {/each}
    <Panel nonInteractive>
        <Header ripple={false}>
            <Button variant="outlined" on:click={newTag}>
                <Label>Add new tag</Label>
                <Icon class="material-icons">add</Icon>
            </Button>
        </Header>
    </Panel>
</Accordion>

<Snackbar bind:this={snackbarTagChangeSaved} timeoutMs=4000>
    <Label>{snackbarTagChangeSavedText}</Label>
</Snackbar>