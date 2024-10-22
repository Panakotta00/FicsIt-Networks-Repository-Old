<script lang="ts">
    import Accordion, {Panel, Header, Content} from '@smui-extra/accordion';
    import {onMount} from "svelte";
    import {getTags, Tag, Tag_Base} from "$lib/api";
    import Textfield from "@smui/textfield";
    import HelperText from "@smui/textfield/helper-text";
    import Button, {Label, Icon} from "@smui/button";
    import Snackbar, {SnackbarComponentDev} from "@smui/snackbar";
    import IconButton, { Icon as ButtonIcon } from '@smui/icon-button';
    import {mutation, operationStore, query} from '@urql/svelte';
    import {noop} from "svelte/internal";

    let tags : Tag_Base[] = []
    let accordion : Accordion
    let panels = {}
    let nameFields = {}
    let snackbarTagChangeSavedText = ""
    let snackbarTagChangeSaved : SnackbarComponentDev
    let tagNegativeID = -1


    const tagsQuery = operationStore(`
        query {
            getAllTags {
                id
                name
                description
            }
        }
    `);

    const deleteTagQuery = mutation({
        query: `
            mutation ($id: ID!) {
                deleteTag(tagId: $id)
            }
        `
    })

    const createTagQuery = mutation({
        query: `
            mutation ($name: String!, $description: String!) {
                createTag(tag: {name: $name, description: $description}) {
                    id
                }
            }
        `
    })

    const updateTagQuery = mutation({
        query: `
            mutation ($id: ID!, $name: String, $description: String) {
                updateTag(tag: {id: $id, name: $name, description: $description})
            }
        `
    })

    query(tagsQuery);

    $: {
        const data = $tagsQuery.data;
        tags = (data && data["getAllTags"] as Tag_Base[]) || []
    }

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

    async function tagChange(tag : Tag_Base) {
        // ignore "New Tag"
        if (tag.name == "New Tag") return;

        var success = false
        if (tag.id < 0) {
            // Create new tag & update tag.id with new DB id or re-fetch all tags
            try {
                const result = await createTagQuery({name: tag.name, description: tag.description})
                console.log(result)
                if (result.data) {
                    tag.id = result.data.id
                    success = true
                }
            } catch (err) {
                console.log(err)
            }
            if (!success) {
                snackbarTagChangeSavedText = `Failed to create Tag '${tag.name}'!`
                snackbarTagChangeSaved.open()
                setTimeout(() => snackbarTagChangeSaved.close(), 2000)
                return
            }
        } else {
            // Update existing tag
            try {
                success = (await updateTagQuery({id: tag.id, name: tag.name, description: tag.description})).data as boolean
            } catch {noop()}
            if (!success) {
                snackbarTagChangeSavedText = `Failed to update Tag '${tag.name}'!`
                snackbarTagChangeSaved.open()
                setTimeout(() => snackbarTagChangeSaved.close(), 2000)
                return
            }
        }

        snackbarTagChangeSavedText = `Tag '${tag.name}' saved!`
        snackbarTagChangeSaved.open()
        setTimeout(() => snackbarTagChangeSaved.close(), 2000)
    }

    async function deleteTag(tag : Tag_Base) {
        if (tag.name != "New Tag") {
            // Remove tag
            var success = false
            try {
                const result = await deleteTagQuery({id: tag.id})
                success = result.data as boolean
            } catch {
                success = false
            }
            if (!success) {
                snackbarTagChangeSavedText = `Failed to remove Tag '${tag.name}'!`
                snackbarTagChangeSaved.open()
                setTimeout(() => snackbarTagChangeSaved.close(), 2000)
                return
            }
        }

        // Remove tag animation
        let panelRemoveAnimation = () => {
            let panel = panels[tag.id].getElement()
            let startHeight = panel.scrollHeight
            panel.classList.add("smui-accordion__panel--removed")
            panel.style.height = startHeight + 'px';
            requestAnimationFrame(function () {
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
        //tags = await getTags() as Tag_Base[]
    })
</script>

{#if $tagsQuery.fetching}
    <h1>Loading tags...</h1>
{:else if tagsQuery.error}
    <h1>Failed to load tags: {$tagsQuery.error.message}</h1>
{:else}
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
{/if}

<Snackbar bind:this={snackbarTagChangeSaved} timeoutMs=4000>
    <Label>{snackbarTagChangeSavedText}</Label>
</Snackbar>