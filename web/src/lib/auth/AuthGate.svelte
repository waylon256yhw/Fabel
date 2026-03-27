<script lang="ts">
  import type { Snippet } from 'svelte'
  import type { AppState } from '../stores/app.svelte.ts'

  let { app, children }: { app: AppState; children: Snippet } = $props()

  let mode: 'login' | 'register' = $state('login')
  let username = $state('')
  let password = $state('')
  let displayName = $state('')
  let error: string | null = $state(null)
  let loading = $state(false)

  async function handleSubmit(e: Event) {
    e.preventDefault()
    if (loading) return
    loading = true
    error = null

    const result =
      mode === 'login'
        ? await app.login(username, password)
        : await app.register(username, password, displayName || undefined)

    if (result) error = result
    loading = false
  }
</script>

{#if !app.authChecked}
  <!-- Loading state -->
  <div class="flex h-screen items-center justify-center bg-background">
    <div class="text-muted-foreground">Loading...</div>
  </div>
{:else if !app.user}
  <!-- Auth form -->
  <div class="flex h-screen items-center justify-center bg-background">
    <div class="w-full max-w-sm space-y-6 p-8">
      <!-- Logo -->
      <div class="text-center">
        <h1 class="text-3xl font-medium tracking-wide text-foreground">Fabel</h1>
        <p class="mt-2 text-sm text-muted-foreground">
          {mode === 'login' ? 'Welcome back' : 'Create your account'}
        </p>
      </div>

      <!-- Form -->
      <form class="space-y-4" onsubmit={handleSubmit}>
        {#if mode === 'register'}
          <div>
            <label for="displayName" class="mb-1 block text-sm text-muted-foreground">Display Name</label>
            <input
              id="displayName"
              type="text"
              bind:value={displayName}
              placeholder="Optional"
              class="w-full rounded-md border border-border bg-card px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground/50 focus:border-primary focus:outline-none focus:ring-1 focus:ring-ring"
            />
          </div>
        {/if}
        <div>
          <label for="username" class="mb-1 block text-sm text-muted-foreground">Username</label>
          <input
            id="username"
            type="text"
            bind:value={username}
            required
            autocomplete="username"
            class="w-full rounded-md border border-border bg-card px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground/50 focus:border-primary focus:outline-none focus:ring-1 focus:ring-ring"
          />
        </div>
        <div>
          <label for="password" class="mb-1 block text-sm text-muted-foreground">Password</label>
          <input
            id="password"
            type="password"
            bind:value={password}
            required
            autocomplete={mode === 'login' ? 'current-password' : 'new-password'}
            class="w-full rounded-md border border-border bg-card px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground/50 focus:border-primary focus:outline-none focus:ring-1 focus:ring-ring"
          />
        </div>

        {#if error}
          <div class="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">{error}</div>
        {/if}

        <button
          type="submit"
          disabled={loading}
          class="w-full rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50"
        >
          {loading ? '...' : mode === 'login' ? 'Sign In' : 'Create Account'}
        </button>
      </form>

      <!-- Toggle -->
      <p class="text-center text-sm text-muted-foreground">
        {#if mode === 'login'}
          No account?
          <button class="text-primary hover:underline" onclick={() => { mode = 'register'; error = null }}>
            Register
          </button>
        {:else}
          Have an account?
          <button class="text-primary hover:underline" onclick={() => { mode = 'login'; error = null }}>
            Sign In
          </button>
        {/if}
      </p>
    </div>
  </div>
{:else}
  {@render children()}
{/if}
