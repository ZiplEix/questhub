import { redirect, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ url, locals: { supabase } }) => {
    const code = url.searchParams.get('code');
    const next = url.searchParams.get('next') ?? '/dashboard';

    if (!code) {
        throw redirect(303, '/login?error=NoCodeProvided');
    }

    // 1. On fait l'échange en dehors du flux de redirection
    let exchangeError = null;
    try {
        const { error } = await supabase.auth.exchangeCodeForSession(code);
        exchangeError = error;
    } catch (err) {
        console.error('Unhandled exception during OAuth code exchange:', err);
        throw redirect(303, '/login?error=ExceptionDuringCallback');
    }

    // 2. On gère le résultat de l'échange
    if (exchangeError) {
        console.error('OAuth callback exchange error:', exchangeError);
        throw redirect(303, '/login?error=OAuthCallbackFailed');
    }

    // 3. Tout est OK, on redirige
    throw redirect(303, next);
};
