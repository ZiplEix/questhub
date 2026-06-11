import { createServerClient } from '@supabase/ssr';
import { type Handle, redirect, isRedirect } from '@sveltejs/kit';
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
    let user = null;
    try {
        const { data } = await event.locals.supabase.auth.getUser();
        user = data?.user || null;
    } catch (err) {
        console.error('Error fetching user in hooks:', err);
    }
    
    if (user) {
        let isBanned = false;
        try {
            // Check if user is banned
            const { data: roleData, error } = await event.locals.supabase
                .from('user_roles')
                .select('is_banned')
                .eq('user_id', user.id)
                .maybeSingle();

            if (error) {
                console.error('Error checking ban status:', error);
            } else {
                isBanned = !!roleData?.is_banned;
            }
        } catch (err) {
            console.error('Failed to query user_roles table:', err);
        }

        if (isBanned) {
            try {
                await event.locals.supabase.auth.signOut();
            } catch (err) {
                console.error('Failed to sign out banned user:', err);
            }
            throw redirect(302, "/login?error=banned");
        }

        event.locals.user = {
            id: user.id,
            email: user.email || '',
            name: user.user_metadata?.name || user.email?.split('@')[0] || 'Joueur',
            image: user.user_metadata?.avatar_url || null
        };

        try {
            // Fetch session
            const { data: { session } } = await event.locals.supabase.auth.getSession();
            event.locals.session = session;
        } catch (err) {
            console.error('Error fetching session in hooks:', err);
            event.locals.session = null;
        }
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
