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
  <div class="flex h-dvh items-center justify-center bg-background">
    <div class="flex flex-col items-center gap-3 text-muted-foreground">
      <!-- Fabel logo spinner -->
      <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 192 192" class="animate-pulse">
        <path d="m163.4 39.45h-11.78v89.49c0 13.12-10.79 23.68-23.18 23.68h-88.75v10.15c0 13.16 10.01 23.76 23.54 23.76h100.2c13.05 0 23.31-11.18 23.31-23.64v-99.8c0-12.46-10.26-23.64-23.31-23.64z" fill="currentColor" opacity="0.3"/>
        <path d="m128.8 5.81h-100.9c-12.62 0-22.25 11.5-22.25 23.31v100.4c0 12.22 10.04 23.13 23.05 23.13h11.02v-91.76c0-11.65 8.2-21.23 20.18-21.23h50.4c5.05 0 7.58 2.23 7.58 9.32v7.17c0 6.79-2.28 10.11-8.61 10.11h-27.42c-7.27 0-9.31 1.68-9.31 9.86v3.32c0 7.46 2.07 9.77 8.33 9.77h20.21c6.32 0 7.57 3.2 7.57 9.44v7.02c0 5.82-2.18 8.13-8.5 8.13h-18.24c-7.14 0-9.3 2.14-9.3 11.54v26.88h56.11c12.89 0 22.93-11.33 22.93-23.38v-99.23c0-12.05-9.93-23.77-22.82-23.77z" fill="currentColor" opacity="0.6"/>
      </svg>
    </div>
  </div>
{:else if !app.user}
  <div class="flex h-dvh items-center justify-center bg-background">
    <div class="w-full max-w-sm rounded-xl border border-border bg-card p-8 shadow-lg" style="box-shadow: 0 4px 24px rgba(0,0,0,0.4);">
      <!-- Logo -->
      <div class="flex flex-col items-center gap-3 pb-6">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 192 192">
          <path d="m163.4 39.45h-11.78v89.49c0 13.12-10.79 23.68-23.18 23.68h-88.75v10.15c0 13.16 10.01 23.76 23.54 23.76h100.2c13.05 0 23.31-11.18 23.31-23.64v-99.8c0-12.46-10.26-23.64-23.31-23.64z" fill="var(--muted-foreground)"/>
          <path d="m128.8 5.81h-100.9c-12.62 0-22.25 11.5-22.25 23.31v100.4c0 12.22 10.04 23.13 23.05 23.13h11.02v-91.76c0-11.65 8.2-21.23 20.18-21.23h50.4c5.05 0 7.58 2.23 7.58 9.32v7.17c0 6.79-2.28 10.11-8.61 10.11h-27.42c-7.27 0-9.31 1.68-9.31 9.86v3.32c0 7.46 2.07 9.77 8.33 9.77h20.21c6.32 0 7.57 3.2 7.57 9.44v7.02c0 5.82-2.18 8.13-8.5 8.13h-18.24c-7.14 0-9.3 2.14-9.3 11.54v26.88h56.11c12.89 0 22.93-11.33 22.93-23.38v-99.23c0-12.05-9.93-23.77-22.82-23.77z" fill="var(--foreground)"/>
        </svg>
        <h1 class="font-story text-2xl font-medium tracking-wide text-foreground">Fabel</h1>
        <p class="text-sm text-muted-foreground">
          {mode === 'login' ? 'Welcome back' : 'Create your account'}
        </p>
      </div>

      <!-- Form -->
      <form class="flex flex-col gap-4" onsubmit={handleSubmit}>
        {#if mode === 'register'}
          <div class="flex flex-col gap-1.5">
            <label for="displayName" class="text-[11px] uppercase tracking-wider text-muted-foreground">Display Name</label>
            <input
              id="displayName"
              type="text"
              bind:value={displayName}
              placeholder="Optional"
              class="w-full rounded-lg border border-border bg-background px-3 py-2.5 text-sm text-foreground placeholder:text-muted-foreground/40 focus:border-primary focus:outline-none focus:ring-1 focus:ring-ring"
            />
          </div>
        {/if}
        <div class="flex flex-col gap-1.5">
          <label for="username" class="text-[11px] uppercase tracking-wider text-muted-foreground">Username</label>
          <input
            id="username"
            type="text"
            bind:value={username}
            required
            autocomplete="username"
            class="w-full rounded-lg border border-border bg-background px-3 py-2.5 text-sm text-foreground placeholder:text-muted-foreground/40 focus:border-primary focus:outline-none focus:ring-1 focus:ring-ring"
          />
        </div>
        <div class="flex flex-col gap-1.5">
          <label for="password" class="text-[11px] uppercase tracking-wider text-muted-foreground">Password</label>
          <input
            id="password"
            type="password"
            bind:value={password}
            required
            autocomplete={mode === 'login' ? 'current-password' : 'new-password'}
            class="w-full rounded-lg border border-border bg-background px-3 py-2.5 text-sm text-foreground placeholder:text-muted-foreground/40 focus:border-primary focus:outline-none focus:ring-1 focus:ring-ring"
          />
        </div>

        {#if error}
          <div class="rounded-lg bg-destructive/10 px-3 py-2 text-sm text-destructive">{error}</div>
        {/if}

        <button
          type="submit"
          disabled={loading}
          class="mt-2 w-full rounded-lg bg-primary px-3 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 hover:shadow-md disabled:opacity-50"
        >
          {loading ? '...' : mode === 'login' ? 'Sign In' : 'Create Account'}
        </button>
      </form>

      <!-- Toggle -->
      <p class="pt-5 text-center text-sm text-muted-foreground">
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
