import { supabase } from './supabaseClient';

export async function sendMessage(message: {
    game_id: string;
    content: string;
    type: string;
    sender_name: string;
    target_id?: string;
}) {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) {
        console.error("No authenticated user available for sendMessage");
        return;
    }

    try {
        let senderName = message.sender_name;

        // Verify/fetch player name if not GM
        const { data: game } = await supabase
            .from('games')
            .select('gm_id')
            .eq('id', message.game_id)
            .single();

        if (game && game.gm_id !== user.id) {
            const { data: gc } = await supabase
                .from('game_characters')
                .select('characters(name)')
                .eq('game_id', message.game_id)
                .eq('user_id', user.id)
                .maybeSingle();

            senderName = (gc?.characters as any)?.name || user.user_metadata?.name || user.email?.split('@')[0] || 'Joueur';
        }

        const { error } = await supabase
            .from('messages')
            .insert({
                game_id: message.game_id,
                sender_id: user.id,
                sender_name: senderName,
                content: message.content,
                type: message.type || 'CHAT_GLOBAL',
                target_id: message.target_id || null
            });

        if (error) throw error;
    } catch (e) {
        console.error("Failed to send message:", e);
    }
}
