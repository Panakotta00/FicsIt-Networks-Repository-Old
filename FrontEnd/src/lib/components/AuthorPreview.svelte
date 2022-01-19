<script lang="ts">
    import {User_Base} from "$lib/api";

    export let author : User_Base = null
</script>

<a href="/user/{author.name}">{author.name}</a>
