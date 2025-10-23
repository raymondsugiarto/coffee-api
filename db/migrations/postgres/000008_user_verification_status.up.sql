ALTER TABLE "user" ADD COLUMN phone_verification_status VARCHAR(100) DEFAULT 'unverified';
ALTER TABLE "user" ADD COLUMN email_verification_status VARCHAR(100) DEFAULT 'unverified';
