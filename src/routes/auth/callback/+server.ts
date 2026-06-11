import { redirect, isRedirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url, locals: { supabase } }) => {
    const code = url.searchParams.get('code');
    const next = url.searchParams.get('next') ?? '/dashboard';

    if (code) {
        try {
            const { error } = await supabase.auth.exchangeCodeForSession(code);
            if (!error) {
                throw redirect(303, next);
            }
            console.error('OAuth callback exchange error:', error);
        } catch (err) {
            if (isRedirect(err)) {
                throw err;
            }
            console.error('Unhandled exception during OAuth code exchange:', err);
        }
    }

    // Redirect to login with error if token exchange fails
    throw redirect(303, '/login?error=OAuthCallbackFailed');
};
