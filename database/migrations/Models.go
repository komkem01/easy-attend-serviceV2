package migrations

import "easy-attend-service/models"

func Models() []any {
	return []any{
		(*models.Gender)(nil),
		(*models.Prefix)(nil),
		(*models.School)(nil),
		(*models.Teacher)(nil),
		(*models.Student)(nil),
		(*models.Classroom)(nil),
		(*models.ClassroomMember)(nil),
		(*models.Attendance)(nil),
		(*models.Log)(nil),
	}
}

func RawBeforeQueryMigrate() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

		// Drop existing tables to avoid timestamp conversion issues
		`DROP TABLE IF EXISTS "attendances" CASCADE;`,
		`DROP TABLE IF EXISTS "logs" CASCADE;`,
		`DROP TABLE IF EXISTS "classroom_members" CASCADE;`,
		`DROP TABLE IF EXISTS "classrooms" CASCADE;`,

		// Convert existing timestamp columns to bigint for teachers table
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'teachers' AND column_name = 'created_at' AND data_type = 'timestamp with time zone') THEN
				ALTER TABLE "teachers" ALTER COLUMN "created_at" TYPE bigint USING EXTRACT(epoch FROM "created_at");
			END IF;
		END $$;`,

		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'teachers' AND column_name = 'updated_at' AND data_type = 'timestamp with time zone') THEN
				ALTER TABLE "teachers" ALTER COLUMN "updated_at" TYPE bigint USING EXTRACT(epoch FROM "updated_at");
			END IF;
		END $$;`,

		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'teachers' AND column_name = 'deleted_at' AND data_type = 'timestamp with time zone') THEN
				ALTER TABLE "teachers" ADD COLUMN "deleted_at_new" bigint;
				UPDATE "teachers" SET "deleted_at_new" = EXTRACT(epoch FROM "deleted_at") WHERE "deleted_at" IS NOT NULL;
				ALTER TABLE "teachers" DROP COLUMN "deleted_at";
				ALTER TABLE "teachers" RENAME COLUMN "deleted_at_new" TO "deleted_at";
			END IF;
		END $$;`,

		// Convert existing timestamp columns to bigint for students table
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'created_at' AND data_type = 'timestamp with time zone') THEN
				ALTER TABLE "students" ALTER COLUMN "created_at" TYPE bigint USING EXTRACT(epoch FROM "created_at");
			END IF;
		END $$;`,

		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'updated_at' AND data_type = 'timestamp with time zone') THEN
				ALTER TABLE "students" ALTER COLUMN "updated_at" TYPE bigint USING EXTRACT(epoch FROM "updated_at");
			END IF;
		END $$;`,

		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'deleted_at' AND data_type = 'timestamp with time zone') THEN
				ALTER TABLE "students" ADD COLUMN "deleted_at_new" bigint;
				UPDATE "students" SET "deleted_at_new" = EXTRACT(epoch FROM "deleted_at") WHERE "deleted_at" IS NOT NULL;
				ALTER TABLE "students" DROP COLUMN "deleted_at";
				ALTER TABLE "students" RENAME COLUMN "deleted_at_new" TO "deleted_at";
			END IF;
		END $$;`,

		// Remove email and phone columns from students if they exist
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'email') THEN
				ALTER TABLE "students" DROP COLUMN "email";
			END IF;
		END $$;`,

		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'phone') THEN
				ALTER TABLE "students" DROP COLUMN "phone";
			END IF;
		END $$;`,

		// Add school_id to teachers table if not exists
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'teachers' AND column_name = 'school_id') THEN
				ALTER TABLE "teachers" ADD COLUMN "school_id" uuid;
			END IF;
		END $$;`,
	}
}

func RawAfterQueryMigrate() []string {
	return []string{
		// Create classrooms table
		`CREATE TABLE IF NOT EXISTS "classrooms" (
			"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			"school_id" uuid NOT NULL,
			"teacher_id" uuid NOT NULL,
			"student_id" uuid NOT NULL,
			"name" varchar(255) NOT NULL,
			"created_at" bigint DEFAULT EXTRACT(epoch FROM NOW()),
			"updated_at" bigint DEFAULT EXTRACT(epoch FROM NOW()),
			"deleted_at" bigint
		);`,

		// Create classroom_members table
		`CREATE TABLE IF NOT EXISTS "classroom_members" (
			"teacher_id" uuid,
			"student_id" uuid,
			"classroom_id" uuid NOT NULL
		);`,

		// Create attendances table
		`CREATE TABLE IF NOT EXISTS "attendances" (
			"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			"classroom_id" uuid NOT NULL,
			"teacher_id" uuid NOT NULL,
			"student_id" uuid NOT NULL,
			"session_date" date NOT NULL,
			"status" varchar(20) NOT NULL,
			"checked_at" bigint NOT NULL,
			"remark" text
		);`,

		// Create logs table
		`CREATE TABLE IF NOT EXISTS "logs" (
			"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			"teacher_id" uuid NOT NULL,
			"action" varchar(100) NOT NULL,
			"detail" text,
			"created_at" bigint DEFAULT EXTRACT(epoch FROM NOW()),
			"school_id" uuid
		);`,

		// Add foreign key constraints
		`ALTER TABLE "classrooms" ADD CONSTRAINT IF NOT EXISTS "fk_classrooms_school" FOREIGN KEY ("school_id") REFERENCES "schools"("id") ON DELETE CASCADE;`,
		`ALTER TABLE "classrooms" ADD CONSTRAINT IF NOT EXISTS "fk_classrooms_teacher" FOREIGN KEY ("teacher_id") REFERENCES "teachers"("id") ON DELETE CASCADE;`,
		`ALTER TABLE "classrooms" ADD CONSTRAINT IF NOT EXISTS "fk_classrooms_student" FOREIGN KEY ("student_id") REFERENCES "students"("id") ON DELETE CASCADE;`,

		`ALTER TABLE "classroom_members" ADD CONSTRAINT IF NOT EXISTS "fk_classroom_members_teacher" FOREIGN KEY ("teacher_id") REFERENCES "teachers"("id") ON DELETE CASCADE;`,
		`ALTER TABLE "classroom_members" ADD CONSTRAINT IF NOT EXISTS "fk_classroom_members_student" FOREIGN KEY ("student_id") REFERENCES "students"("id") ON DELETE CASCADE;`,
		`ALTER TABLE "classroom_members" ADD CONSTRAINT IF NOT EXISTS "fk_classroom_members_classroom" FOREIGN KEY ("classroom_id") REFERENCES "classrooms"("id") ON DELETE CASCADE;`,

		`ALTER TABLE "attendances" ADD CONSTRAINT IF NOT EXISTS "fk_attendances_classroom" FOREIGN KEY ("classroom_id") REFERENCES "classrooms"("id") ON DELETE CASCADE;`,
		`ALTER TABLE "attendances" ADD CONSTRAINT IF NOT EXISTS "fk_attendances_teacher" FOREIGN KEY ("teacher_id") REFERENCES "teachers"("id") ON DELETE CASCADE;`,
		`ALTER TABLE "attendances" ADD CONSTRAINT IF NOT EXISTS "fk_attendances_student" FOREIGN KEY ("student_id") REFERENCES "students"("id") ON DELETE CASCADE;`,

		`ALTER TABLE "logs" ADD CONSTRAINT IF NOT EXISTS "fk_logs_teacher" FOREIGN KEY ("teacher_id") REFERENCES "teachers"("id") ON DELETE CASCADE;`,
		`ALTER TABLE "logs" ADD CONSTRAINT IF NOT EXISTS "fk_logs_school" FOREIGN KEY ("school_id") REFERENCES "schools"("id") ON DELETE CASCADE;`,

		// Add indexes
		`CREATE INDEX IF NOT EXISTS "idx_classrooms_school_id" ON "classrooms"("school_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_classrooms_teacher_id" ON "classrooms"("teacher_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_classrooms_student_id" ON "classrooms"("student_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_classrooms_deleted_at" ON "classrooms"("deleted_at");`,

		`CREATE INDEX IF NOT EXISTS "idx_classroom_members_teacher_id" ON "classroom_members"("teacher_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_classroom_members_student_id" ON "classroom_members"("student_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_classroom_members_classroom_id" ON "classroom_members"("classroom_id");`,

		`CREATE INDEX IF NOT EXISTS "idx_attendances_classroom_id" ON "attendances"("classroom_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_attendances_teacher_id" ON "attendances"("teacher_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_attendances_student_id" ON "attendances"("student_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_attendances_session_date" ON "attendances"("session_date");`,
		// Composite indexes for common queries
		`CREATE INDEX IF NOT EXISTS "idx_attendances_classroom_student_date" ON "attendances"("classroom_id", "student_id", "session_date");`,
		`CREATE INDEX IF NOT EXISTS "idx_attendances_student_date" ON "attendances"("student_id", "session_date");`,
		`CREATE INDEX IF NOT EXISTS "idx_attendances_classroom_date" ON "attendances"("classroom_id", "session_date");`,

		`CREATE INDEX IF NOT EXISTS "idx_logs_teacher_id" ON "logs"("teacher_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_logs_school_id" ON "logs"("school_id");`,
		`CREATE INDEX IF NOT EXISTS "idx_logs_created_at" ON "logs"("created_at");`,
	}
}
