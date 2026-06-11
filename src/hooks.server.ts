import { createServerClient } from '@supabase/ssr';
import { type Handle, redirect } from '@sveltejs/kit';
import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public';

const protectedRoutes = [
    "/dashboard",
    "/table",
    "/reset-password"
];

export const handle: Handle = async ({ event, resolve }) => {
    // 1. Initialize Supabase Server Client
    event.locals.supabase = createServerClient(
        PUBLIC_SUPABASE_URL,
        PUBLIC_SUPABASE_ANON_KEY,
        {
            cookies: {
                getAll() {
                    return event.cookies.getAll();
                },
                setAll(cookiesToSet: { name: string; value: string; options: any }[]) {
                    cookiesToSet.forEach(({ name, value, options }) => {
                        event.cookies.set(name, value, { ...options, path: '/' });
                    });
                },
            },
        }
    );

    // 2. Fetch authenticated user (securely validates JWT against Supabase server)
    const { data: { user } } = await event.locals.supabase.auth.getUser();
    
    if (user) {
        // Check if user is banned
        const { data: roleData } = await event.locals.supabase
            .from('user_roles')
            .select('is_banned')
            .eq('user_id', user.id)
            .maybeSingle();

        if (roleData?.is_banned) {
            await event.locals.supabase.auth.signOut();
            throw redirect(302, "/login?error=banned");
        }

        event.locals.user = {
            id: user.id,
            email: user.email || '',
            name: user.user_metadata?.name || user.email?.split('@')[0] || 'Joueur',
            image: user.user_metadata?.avatar_url || null
        };
        // Fetch session
        const { data: { session } } = await event.locals.supabase.auth.getSession();
        event.locals.session = session;
    } else {
        event.locals.user = null;
        event.locals.session = null;
    }

    // 3. Handle Route Guards
    const pathname = event.url.pathname;
    const isProtected = protectedRoutes.some((route) => pathname.startsWith(route));

    if (isProtected && !event.locals.user) {
        throw redirect(302, "/login");
    }

    // 4. Resolve request, allowing Supabase headers
    return resolve(event, {
        filterSerializedResponseHeaders(name) {
            return name === 'content-range' || name === 'x-supabase-api-version';
        },
    });
};
