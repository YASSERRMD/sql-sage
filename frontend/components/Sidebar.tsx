"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useTheme } from "next-themes";
import { useMe, useLogout } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";

const NAV = [
  { href: "/", label: "Dashboard" },
  { href: "/workspace", label: "Workspace" },
  { href: "/history", label: "History" },
  { href: "/providers", label: "Providers" },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const { data: user } = useMe();
  const logout = useLogout();
  const { theme, setTheme } = useTheme();

  return (
    <aside className="hidden w-60 shrink-0 border-r bg-card/50 p-4 md:block">
      <div className="mb-6 flex items-center gap-2">
        <div className="grid h-8 w-8 place-items-center rounded-md bg-primary text-primary-foreground">
          S
        </div>
        <div>
          <p className="text-sm font-semibold">SQL-Sage</p>
          <p className="text-xs text-muted-foreground">Static analysis</p>
        </div>
      </div>
      <nav className="space-y-1">
        {NAV.map((n) => {
          const active = pathname === n.href;
          return (
            <Link
              key={n.href}
              href={n.href}
              className={`block rounded-md px-3 py-2 text-sm transition-colors ${
                active
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-accent hover:text-accent-foreground"
              }`}
            >
              {n.label}
            </Link>
          );
        })}
      </nav>
      <div className="mt-6 border-t pt-4 text-xs text-muted-foreground">
        {user && (
          <p className="mb-2">
            Signed in as
            <br />
            <span className="font-medium text-foreground">{user.name}</span>
          </p>
        )}
        <Button
          variant="outline"
          size="sm"
          className="w-full"
          onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
        >
          {theme === "dark" ? "Light" : "Dark"} mode
        </Button>
        <Button
          variant="ghost"
          size="sm"
          className="mt-2 w-full"
          onClick={() =>
            logout.mutate(undefined, {
              onSuccess: () => router.push("/login"),
            })
          }
        >
          Sign out
        </Button>
      </div>
    </aside>
  );
}
