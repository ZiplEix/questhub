import { writable } from 'svelte/store';
import { supabase } from './supabaseClient';

export interface UserSession {
    data: {
        session: any;
        user: {
            id: string;
            name: string;
            email: string;
            image: string | null;
            createdAt: Date;
            updatedAt: Date;
            emailVerified: boolean;
        };
    } | null;
    isPending: boolean;
    error: any;
}

const sessionStore = writable<UserSession>({
    data: null,
    isPending: true,
    error: null
});

// Helper to construct Svelte-compatible Session object conforming to frontend expectations
function mapSession(session: any) {
    if (!session) return null;
    return {
        session: session,
        user: {
            id: session.user.id,
            name: session.user.user_metadata?.name || session.user.email?.split('@')[0] || 'Joueur',
            email: session.user.email || '',
            image: session.user.user_metadata?.avatar_url || null,
            createdAt: new Date(session.user.created_at || Date.now()),
            updatedAt: new Date(session.user.updated_at || Date.now()),
            emailVerified: !!session.user.email_confirmed_at
        }
    };
}

// Setup active listener to track user state and metadata
supabase.auth.getSession().then(({ data: { session } }) => {
    sessionStore.set({
        data: mapSession(session),
        isPending: false,
        error: null
    });
});

supabase.auth.onAuthStateChange((event, session) => {
    sessionStore.set({
        data: mapSession(session),
        isPending: false,
        error: null
    });
});

export const authClient = {
    // Getter for direct access if needed
    get session() {
        let current: UserSession | null = null;
        sessionStore.subscribe(val => { current = val; })();
        return current;
    },

    useSession() {
        return {
            subscribe: sessionStore.subscribe
        };
    },

    // Mock token function returning a mock or real JWT token for retro-compatibility
    async token() {
        const { data: { session } } = await supabase.auth.getSession();
        return {
            data: {
                token: session?.access_token || null
            },
            error: null
        };
    },

    // Mock getSession function
    async getSession() {
        const { data: { session } } = await supabase.auth.getSession();
        return {
            data: mapSession(session),
            error: null
        };
    },

    signUp: {
        async email(
            data: { email: string; password: string; name: string; callbackURL?: string },
            options?: {
                onError?: (ctx: { error: { message: string } }) => void;
                onResponse?: () => void;
                onSuccess?: (ctx: { data: any }) => void;
            }
        ) {
            try {
                const { data: resData, error } = await supabase.auth.signUp({
                    email: data.email,
                    password: data.password,
                    options: {
                        data: { name: data.name },
                        emailRedirectTo: data.callbackURL ? (window.location.origin + data.callbackURL) : undefined
                    }
                });

                if (error) throw error;

                options?.onSuccess?.({ data: resData });
                return resData;
            } catch (err: any) {
                options?.onError?.({ error: { message: err.message || 'Erreur lors de l\'inscription.' } });
                throw err;
            } finally {
                options?.onResponse?.();
            }
        }
    },

    signIn: {
        async email(
            data: { email: string; password: string; callbackURL?: string },
            options?: {
                onError?: (ctx: { error: { message: string } }) => void;
                onResponse?: () => void;
                onSuccess?: (ctx: { data: any }) => void;
            }
        ) {
            try {
                const { data: resData, error } = await supabase.auth.signInWithPassword({
                    email: data.email,
                    password: data.password
                });

                if (error) throw error;

                options?.onSuccess?.({ data: resData });

                if (data.callbackURL) {
                    window.location.href = data.callbackURL;
                }
                return resData;
            } catch (err: any) {
                options?.onError?.({ error: { message: err.message || 'Identifiants invalides.' } });
                throw err;
            } finally {
                options?.onResponse?.();
            }
        },

        async social(
            data: { provider: 'google'; callbackURL: string; errorCallbackURL?: string },
            options?: {
                onError?: (ctx: { error: { message: string } }) => void;
                onResponse?: () => void;
                onSuccess?: (ctx: { data: any }) => void;
            }
        ) {
            try {
                const { data: resData, error } = await supabase.auth.signInWithOAuth({
                    provider: data.provider,
                    options: {
                        redirectTo: window.location.origin + '/auth/callback?next=' + encodeURIComponent(data.callbackURL)
                    }
                });

                if (error) throw error;

                if (resData?.url) {
                    window.location.href = resData.url;
                }

                options?.onSuccess?.({ data: resData });
                return resData;
            } catch (err: any) {
                options?.onError?.({ error: { message: err.message || 'Erreur lors de la connexion sociale.' } });
                throw err;
            } finally {
                options?.onResponse?.();
            }
        }
    },

    async signOut(options?: {
        fetchOptions?: {
            onError?: (ctx: { error: { message: string } }) => void;
            onSuccess?: () => void;
        }
    }) {
        try {
            const { error } = await supabase.auth.signOut();
            if (error) throw error;
            options?.fetchOptions?.onSuccess?.();
        } catch (err: any) {
            options?.fetchOptions?.onError?.({ error: { message: err.message || 'Erreur lors de la déconnexion.' } });
            throw err;
        }
    }
};
