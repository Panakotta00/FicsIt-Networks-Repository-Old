<script lang="ts">
    import PackagePreview from "../lib/components/PackagePreview.svelte";
    import {listPackages, Package} from "$lib/api.ts";

    async function getData() : Promise<Package[]> {
        return await listPackages(0, 50)
    }
</script>

<h1>Welcome to SvelteKit</h1>

{#await listPackages()}
    <h1>Loading...</h1>
{:then packages}
    {#each packages as pack}
        <a href="/package/{pack.name}">
            <PackagePreview packageStruct={pack} />
        </a>
    {/each}
{:catch _}
    <h1>Failed to fetch packages!</h1>
{/await}