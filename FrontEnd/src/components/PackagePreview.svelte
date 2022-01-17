<script lang="ts">
    import {getPackage, listPackages, Package} from "$lib/api.ts";
    import {onMount} from "svelte";

    export let packageId : null|bigint = null
    export let packageStruct : null|Package = null

    async function getData() : Promise<Package> {
        if (packageStruct) {
            return packageStruct
        } else {
            return await getPackage(packageId)
        }
    }
</script>

{#await getData()}
    <h1>Loading...</h1>
{:then data}
    <h1>{data.displayName}</h1>
    <p>{data.description}</p>
{:catch err}
    <h1>Failed to Load!</h1>
{/await}