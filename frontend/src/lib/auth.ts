import { Pool } from 'pg';
import { env } from '$env/dynamic/private';
import { betterAuth } from 'better-auth';
import { jwt } from 'better-auth/plugins';

const pool = new Pool({
    connectionString: env.DATABASE_URL
});

export const auth = betterAuth({
    secret: env.BETTER_AUTH_SECRET,
    url: env.BETTER_AUTH_URL,
    database: pool,
    emailAndPassword: { enabled: true },
    trustedOrigins: env.BETTER_AUTH_TRUSTED_ORIGINS ? env.BETTER_AUTH_TRUSTED_ORIGINS.split(',') : [],

    plugins: [
        jwt(),
    ],

    socialProviders: {
        google: {
            clientId: env.GOOGLE_CLIENT_ID,
            clientSecret: env.GOOGLE_CLIENT_SECRET,
            prompt: "select_account"
        }
    }
})
