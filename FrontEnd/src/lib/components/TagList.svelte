<script lang="ts">
    import {getTags, Tag_Base} from "$lib/api";
    import {onMount} from "svelte";
    import Chip, { Set, TrailingAction, Text } from '@smui/chips';
    import {placeCaretAtEnd, setTextRange} from "$lib/util";
    import MenuSurface, { MenuSurfaceComponentDev } from '@smui/menu-surface';

    export let tags : Tag_Base[] = null
    export let editable : boolean|null = undefined
    export let shake = false
    export let filteredAvailableTags : Tag_Base[] = null
    export let filteredUnavailableTags : Tag_Base[] = null

    export var allowedTags : Tag_Base[]

    let newTagText : string

    function setTagText(text : string) {
        newTagText = text
        newTag.textContent = newTagText
    }

    function filterAvailableTags(tagList : Tag_Base[], currentTags : Tag_Base[], filterText : string) : [Tag_Base[], Tag_Base[]] {
        if (!tagList || !currentTags) return [tagList, tagList]
        var unfiltered = tagList.filter(tag => !currentTags.find(t => t.id == tag.id))
        var filtered = unfiltered.filter(tag => {
            return !newTag || tag.name.startsWith(filterText)
        })
        unfiltered = unfiltered.filter(tag => filtered.findIndex(t => t.id === tag.id) === -1)
        return [filtered, unfiltered]
    }

    $: {
        [filteredAvailableTags, filteredUnavailableTags] = filterAvailableTags(allowedTags, tags, newTagText)
    }

    let newTag : HTMLInputElement = null;
    let newTagContainer : HTMLInputElement = null;
    let surface: MenuSurfaceComponentDev;

    async function removeTag(tagId : bigint) {
        tags = tags.filter(tag => tag.id != tagId)
        tags = tags
    }

    function addTag(newTag : string|bigint|Tag_Base) {
        if (!allowedTags) return false
        let tagToAdd = allowedTags.find(tag => {
            if (typeof newTag == "bigint") {
                return newTag == tag.id
            } else if (typeof newTag == "string") {
                return newTag == tag.name || newTag == tag.id
            } else if (newTag instanceof Tag_Base) {
                return newTag.id == tag.id
            }
            return false
        }) as Tag_Base
        if (tagToAdd && !tags.find(tag => tag.id == tagToAdd.id)) {
            tags.push(tagToAdd)
            tags = tags
            return true
        }
        return false
    }

    function newTagKeydown(e : KeyboardEvent) {
        if (e.code == "Backspace") {
            if (newTag.textContent == "") {
                setTagText(tags.pop().name)
                placeCaretAtEnd(newTag)
                tags = tags
                e.preventDefault()
            }
        } else if (e.code == "Enter") {
            e.preventDefault()
            if (addTag(newTag.textContent)) {
                setTagText("")
            } else {
                shake = true
                setTimeout(() => shake = false, 500)
            }
        } else {
            const newText = newTagText + e.key
            const [available] = filterAvailableTags(allowedTags, tags, newText)
            if (available && available.length > 0) {
                newTag.textContent = available[0].name
                setTextRange(newTag, newTagText.length + 1, newTag.textContent.length)
                e.preventDefault()
                newTagText = newText
            }
        }
    }

    function getFocus() {
        if (newTag && !newTag.contains(document.activeElement)) {
            newTag.focus()
            placeCaretAtEnd(newTag)
        }
    }

    onMount(async () => {
        allowedTags = await getTags()
    })
</script>

<div class="tags" on:click={getFocus}>
    <Set class="tagList" chips={tags} let:chip key={tag => tag.name} nonInteractive>
        <!--Wrapper-->
            <Chip {chip} shouldRemoveOnTrailingIconClick=true on:SMUIChip:removal={() => {tags = tags; setTimeout(() => getFocus(), 100)}}>
                <Text>{chip.name}</Text>
                {#if editable}
                    <TrailingAction icon$class="material-icons">cancel</TrailingAction>
                {/if}
            </Chip>
            <!--Tooltip xPos="start">{chip.description}</Tooltip>
        </Wrapper-->
    </Set>
    {#if editable}
        <div id="newTagContainer" bind:this={newTagContainer} on:focusin={() => surface.setOpen(true)} on:focusout={() => setTimeout(() => {if (!newTagContainer.contains(document.activeElement)) surface.setOpen(false)}, 200)}>
            <MenuSurface bind:this={surface} managed="true" anchorCorner="BOTTOM_LEFT" anchorElement={newTag}>
                <div style="margin: 1rem">
                    <h1>Available Tags</h1>
                    <div class="flex flex-wrap m-1">
                        <Set chips={filteredAvailableTags} let:chip key={tag => tag.name}>
                            <!--Wrapper-->
                                <Chip {chip} on:SMUIChip:interaction={() => addTag(chip.name)}>
                                    <Text>{chip.name}</Text>
                                </Chip>
                                <!--Tooltip xPos="start">{chip.description}</Tooltip>
                            </Wrapper-->
                        </Set>
                    </div>
                    <div class="flex flex-wrap m-1">
                        <Set chips={filteredUnavailableTags} let:chip key={tag => tag.name}>
                            <!--Wrapper-->
                                <Chip {chip} on:SMUIChip:interaction={() => addTag(chip.name)}>
                                    <Text>{chip.name}</Text>
                                </Chip>
                                <!--Tooltip xPos="start">{chip.description}</Tooltip>
                            </Wrapper-->
                        </Set>
                    </div>
                </div>
            </MenuSurface>
            <div id="newTagScroll">
                <span type="text" id="newTag" spellcheck="false" contenteditable="true"  role="textbox" class:shake bind:this={newTag} on:keydown={newTagKeydown} on:input={() => newTagText = newTag.textContent}></span>
            </div>
        </div>
    {/if}
</div>

<style>
    .tags {
        display: flex;
        flex-wrap: wrap;
    }

    .tagList {
        background-color: blue;
        flex: auto;
    }

    #newTagContainer {
        display: flex;
        flex: 1;
        max-width: 100%;
        width: 100%;
    }

    #newTagScroll {
        display: flex;
        min-width: 9rem;
        width: 100%;
    }

    #newTag {
        @apply p-3 mt-auto mb-auto;
        overflow: auto;
        border: none;
        outline: none;
        background: transparent;
        width: 100%;
    }
</style>