// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: {
				id: string;
				email: string;
				name: string | null;
				image?: string | null;
			} | null;
			session: unknown | null;
			supabase: import('@supabase/supabase-js').SupabaseClient;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export { };
