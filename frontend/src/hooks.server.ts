import { building } from "$app/environment";
import { auth } from "$lib/auth";
import { redirect } from "@sveltejs/kit";
import { svelteKitHandler } from "better-auth/svelte-kit";

const protectedRoutes = [
    "/dashboard",
    "/table"
];

export const handle = async ({ event, resolve }) => {
    const session = await auth.api.getSession({
        headers: event.request.headers,
    });

    event.locals.session = session?.session ?? null;
    event.locals.user = session?.user ?? null;

    const pathname = event.url.pathname;

    const isProtected = protectedRoutes.some((route) => pathname.startsWith(route));

    if (isProtected && !event.locals.user) {
        throw redirect(302, "/login");
    }

    return svelteKitHandler({ event, resolve, auth, building });
}
