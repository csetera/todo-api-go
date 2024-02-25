CREATE TABLE public.to_do_item_entities (
	id bigserial NOT NULL,
	description text NULL,
	completed bool NULL,
	due_date timestamptz NULL,
	completed_at timestamptz NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT to_do_item_entities_pkey PRIMARY KEY (id)
);
