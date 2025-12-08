import { api } from "$lib/api";
import { authClient } from "$lib/auth-client";

export async function sendMessage(message: {
    game_id: string;
    content: string;
    type: string;
    sender_name: string;
    target_id?: string;
}) {
    const { data } = await authClient.token();
    const token = data?.token;

    if (!token) {
        console.error("No auth token available for sendMessage");
        return;
    }
    if (!message.game_id) {
        console.error("No game_id provided in message");
        return;
    }

    try {
        await api.post(`/table/${message.game_id}/chat`, message, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
    } catch (e) {
        console.error("Failed to send message:", e);
    }
}
