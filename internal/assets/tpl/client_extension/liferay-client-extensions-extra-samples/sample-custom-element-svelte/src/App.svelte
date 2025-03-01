<svelte:options customElement={{
  tag: 'sample-custom-element-svelte',
  shadow: 'none'
}}/>

<script>
  import LiferayService from "./lib/LiferayService.js";

  let user = $state({name: null})
  let error = $state(false)
  let errorMessage = $state('')
  let loading = $state(true)

  LiferayService.get('/o/headless-admin-user/v1.0/my-user-account')
  .then((response) => {
    user = response
    loading = false
  })
  .catch((err) => {
    error = true
    errorMessage = err.message
  })
  .finally(() => {
    loading = false
  })
</script>

<div>
  {#if loading}
    <span
        aria-hidden="true"
        class="loading-animation-squares loading-animation-primary loading-animation-md"
    ></span>
  {:else if error}
    <div class="alert alert-danger">
      <strong class="lead">Error:</strong>{errorMessage}
    </div>
  {:else}
    <h1 style="color: var(--primary)">
      Hello from Svelte{', ' + user.name + '!' || '!'}
    </h1>
  {/if}
</div>