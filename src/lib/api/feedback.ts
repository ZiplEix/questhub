import { supabase } from '../supabaseClient';

export interface SupportTicketPayload {
    email: string;
    category: 'CONTACT' | 'BUG' | 'RECLAMATION';
    subject: string;
    message: string;
    attachmentFile?: File;
}

export async function submitSupportTicket(payload: SupportTicketPayload): Promise<void> {
    let attachmentUrl: string | null = null;

    // 1. Upload attachment if present
    if (payload.attachmentFile) {
        const file = payload.attachmentFile;
        const fileExt = file.name.split('.').pop();
        const randName = Math.random().toString(36).substring(2, 15);
        const fileName = `${randName}-${Date.now()}.${fileExt}`;
        const filePath = `feedback/${fileName}`;

        const { error: uploadError } = await supabase.storage
            .from('images')
            .upload(filePath, file);

        if (uploadError) {
            console.error('Failed to upload feedback attachment:', uploadError);
            throw new Error('Erreur lors du téléchargement de la pièce jointe.');
        }

        // Retrieve public URL
        const { data } = supabase.storage
            .from('images')
            .getPublicUrl(filePath);

        attachmentUrl = data?.publicUrl || null;
    }

    // 2. Resolve optional logged-in user
    const { data: { user } } = await supabase.auth.getUser();
    const userId = user ? user.id : null;

    // 3. Insert ticket into support_tickets table
    const { error: insertError } = await supabase
        .from('support_tickets')
        .insert({
            user_id: userId,
            email: payload.email,
            category: payload.category,
            subject: payload.subject,
            message: payload.message,
            attachment_url: attachmentUrl,
            status: 'PENDING'
        });

    if (insertError) {
        console.error('Failed to insert support ticket:', insertError);
        throw insertError;
    }
}
