<script lang="ts">
    import {getTags, Tag_Base} from "$lib/api";
    import {onMount} from "svelte";

    export let tags : Tag_Base[] = null
    export let editable : boolean|null = undefined
    export let shake = false

    export var allowedTags : Tag_Base[]

    let newTag : HTMLInputElement = null;

    async function removeTag(tagId : bigint) {
        tags = tags.filter(tag => tag.id != tagId)
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
        console.log(e.code)
        if (e.code == "Backspace") {
            if (newTag.value == "") {
                newTag.value = tags.pop().name
                tags = tags
                e.preventDefault()
            }
        } else if (e.code == "Enter") {
            if (addTag(newTag.value)) {
                newTag.value = ""
            } else {
                shake = true
                setTimeout(() => shake = false, 500)
            }
        }
    }

    onMount(async () => {
        allowedTags = await getTags()
    })
</script>

<div class="tags" on:click={newTag.focus()}>
    {#each tags as tag}
        <div class="tag">
            <p>{tag.name}</p>
            {#if editable}
                <button class="removeTag" on:click={removeTag(tag.id)}><p>×</p></button>
            {/if}
        </div>
    {/each}
    {#if editable}
        <input type="text" id="newTag" spellcheck="false" class:shake bind:this={newTag} on:keydown={newTagKeydown} />
    {/if}
</div>

<style>
    .tags {
        cursor: default;
        display: flex;
        flex-wrap: wrap;
    }

    .tags > * {
        flex: auto auto 1;
    }

    .tag {
        @apply bg-gray-500 p-0.5 pl-2 pr-2 m-0.5 text-white;
        border-radius: 1rem;
        display: inline-flex;
    }

    .tag > p {
        @apply p-0 pr-1;
    }

    .removeTag {
        @apply m-0 m-auto;
        border-radius: 50%;
        width: 1rem;
        height: 1rem;
        line-height: 1rem;
    }

    .removeTag:hover {
        @apply bg-gray-700;
    }

    .removeTag > p {
        position: relative;
        top: 0.5px;
    }

    #newTag {
        @apply p-1;
        border: none;
        outline: none;
        min-width: 3rem;
        flex: 1;
    }
</style>