import { supabase } from '../supabaseClient';

export interface AdminStats {
    total_users: number;
    total_games: number;
    total_templates: number;
    pending_tickets: number;
    banned_users: number;
}

export interface AdminUser {
    id: string;
    email: string;
    name: string;
    created_at: string;
    role: 'admin' | 'moderator' | null;
    is_banned: boolean;
}

export interface AdminGame {
    id: string;
    name: string;
    gm_name: string;
    gm_email: string;
    state: string;
    is_active: boolean;
    created_at: string;
    player_count: number;
}

export interface AdminTemplate {
    id: string;
    name: string;
    type: string;
    creator_name: string;
    is_public: boolean;
    uses: number;
    created_at: string;
    description: string | null;
    data: any;
}

export interface AdminTicket {
    id: string;
    created_at: string;
    email: string;
    category: 'CONTACT' | 'BUG' | 'RECLAMATION';
    subject: string;
    message: string;
    attachment_url: string | null;
    status: 'PENDING' | 'IN_PROGRESS' | 'RESOLVED' | 'CLOSED';
    user_name: string;
}

export async function fetchCurrentUserRole(): Promise<{ role: 'admin' | 'moderator' | null; is_banned: boolean }> {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) return { role: null, is_banned: false };

    const { data, error } = await supabase
        .from('user_roles')
        .select('role, is_banned')
        .eq('user_id', user.id)
        .maybeSingle();

    if (error) {
        console.error("Failed to fetch current user role", error);
        return { role: null, is_banned: false };
    }

    return data || { role: null, is_banned: false };
}

export async function fetchAdminStats(): Promise<AdminStats> {
    const { data, error } = await supabase.rpc('admin_get_stats');
    if (error) throw error;
    return data;
}

export async function fetchAdminUsers(): Promise<AdminUser[]> {
    const { data, error } = await supabase.rpc('admin_get_users');
    if (error) throw error;
    return data || [];
}

export async function setAdminUserRole(targetUserId: string, role: 'admin' | 'moderator' | null): Promise<void> {
    const { error } = await supabase.rpc('admin_set_user_role', {
        p_target_user_id: targetUserId,
        p_role: role
    });
    if (error) throw error;
}

export async function setAdminUserBanned(targetUserId: string, isBanned: boolean): Promise<void> {
    const { error } = await supabase.rpc('admin_set_user_banned', {
        p_target_user_id: targetUserId,
        p_is_banned: isBanned
    });
    if (error) throw error;
}

export async function fetchAdminGames(): Promise<AdminGame[]> {
    const { data, error } = await supabase.rpc('admin_get_games');
    if (error) throw error;
    return data || [];
}

export async function archiveAdminGame(gameId: string, isActive: boolean): Promise<void> {
    const { error } = await supabase.rpc('admin_archive_game', {
        p_game_id: gameId,
        p_is_active: isActive
    });
    if (error) throw error;
}

export async function deleteAdminGame(gameId: string): Promise<void> {
    const { error } = await supabase.rpc('admin_delete_game', {
        p_game_id: gameId
    });
    if (error) throw error;
}

export async function fetchAdminTemplates(): Promise<AdminTemplate[]> {
    const { data, error } = await supabase.rpc('admin_get_templates');
    if (error) throw error;
    return data || [];
}

export async function toggleAdminTemplatePublic(templateId: string, isPublic: boolean): Promise<void> {
    const { error } = await supabase.rpc('admin_toggle_template_public', {
        p_template_id: templateId,
        p_is_public: isPublic
    });
    if (error) throw error;
}

export async function deleteAdminTemplate(templateId: string): Promise<void> {
    const { error } = await supabase.rpc('admin_delete_template', {
        p_template_id: templateId
    });
    if (error) throw error;
}

export async function fetchAdminTickets(): Promise<AdminTicket[]> {
    const { data, error } = await supabase.rpc('admin_get_tickets');
    if (error) throw error;
    return data || [];
}

export async function updateAdminTicketStatus(ticketId: string, status: string): Promise<void> {
    const { error } = await supabase.rpc('admin_update_ticket_status', {
        p_ticket_id: ticketId,
        p_status: status
    });
    if (error) throw error;
}

export interface AdminUserDetails {
    games_mastered: number;
    games_played: number;
    characters_created: number;
    games: Array<{
        id: string;
        name: string;
        image_url: string | null;
        is_gm: boolean;
        joined_at: string;
    }>;
}

export interface AdminGameDetails {
    boards_count: number;
    messages_count: number;
    characters_count: number;
    players: Array<{
        user_id: string;
        name: string;
        email: string;
        joined_at: string;
        ping_color: string;
        character: {
            name: string;
            avatar_url: string | null;
        } | null;
    }>;
}

export interface AdminTemplateDetails {
    data: any;
    creator_name: string;
}

export async function fetchAdminUserDetails(userId: string): Promise<AdminUserDetails> {
    const { data, error } = await supabase.rpc('admin_get_user_details', { p_user_id: userId });
    if (error) throw error;
    return data;
}

export async function fetchAdminGameDetails(gameId: string): Promise<AdminGameDetails> {
    const { data, error } = await supabase.rpc('admin_get_game_details', { p_game_id: gameId });
    if (error) throw error;
    return data;
}

export async function fetchAdminTemplateDetails(templateId: string): Promise<AdminTemplateDetails> {
    const { data, error } = await supabase.rpc('admin_get_template_details', { p_template_id: templateId });
    if (error) throw error;
    return data;
}
