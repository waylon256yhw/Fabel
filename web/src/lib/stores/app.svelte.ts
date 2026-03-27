export type WorkspaceId = "chat" | "characters" | "worldbooks" | "search";

const VALID_WORKSPACES: WorkspaceId[] = ["chat", "characters", "worldbooks", "search"];

interface AuthUser {
  id: string;
  username: string;
  display_name: string;
  role: "admin" | "user";
}

export class AppState {
  // Auth
  user: AuthUser | null = $state(null);
  authChecked = $state(false);

  // Shell UI
  activeWorkspace: WorkspaceId = $state("chat");
  sidebarOpen = $state(true);
  inspectorOpen = $state(false);

  // TopBar center — set by active workspace
  topBarTitle = $state("");
  topBarSubtitle = $state("");

  constructor() {
    this.syncFromHash();
    window.addEventListener("hashchange", () => this.syncFromHash());
  }

  private syncFromHash() {
    const hash = window.location.hash.slice(1) || "chat";
    this.activeWorkspace = VALID_WORKSPACES.includes(hash as WorkspaceId)
      ? (hash as WorkspaceId)
      : "chat";
  }

  navigate(ws: WorkspaceId) {
    window.location.hash = ws;
  }

  toggleSidebar() {
    this.sidebarOpen = !this.sidebarOpen;
  }

  toggleInspector() {
    this.inspectorOpen = !this.inspectorOpen;
  }

  // ── Auth ──

  async checkAuth() {
    try {
      const res = await fetch("/api/auth/me");
      if (res.ok) this.user = await res.json();
    } catch {
      // not logged in
    }
    this.authChecked = true;
  }

  async login(username: string, password: string): Promise<string | null> {
    const res = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ message: "Login failed" }));
      return err.message;
    }
    this.user = await res.json();
    return null;
  }

  async register(username: string, password: string, displayName?: string): Promise<string | null> {
    const res = await fetch("/api/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password, display_name: displayName }),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ message: "Registration failed" }));
      return err.message;
    }
    this.user = await res.json();
    return null;
  }

  async logout() {
    await fetch("/api/auth/logout", { method: "POST" });
    this.user = null;
  }
}
