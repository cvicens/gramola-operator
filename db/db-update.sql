CREATE OR REPLACE FUNCTION db_update(text) RETURNS integer
LANGUAGE plpgsql
AS $$
DECLARE
    
    -- Declare aliases for user input.
    user_name ALIAS FOR $1;
    
    -- Declare a variables
    found_table TEXT;
    found_sequence TEXT;
BEGIN
    SELECT table_name INTO found_table FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'event';
    IF NOT FOUND THEN
        --
        -- Name: event; Type: TABLE; Schema: public; Owner: {{DB_USERNAME}}
        --

        CREATE TABLE IF NOT EXISTS public.event (
            id bigint NOT NULL,
            address character varying(255),
            artist character varying(255),
            city character varying(255),
            country character varying(255),
            date character varying(255),
            description character varying(255),
            end_time character varying(255),
            image character varying(255),
            location character varying(255),
            name character varying(255),
            province character varying(255),
            start_time character varying(255)
        );

        EXECUTE 'ALTER TABLE public.event OWNER TO ' || user_name;
        -- ALTER TABLE public.event OWNER TO user_name;

        --
        -- Name: event event_pkey; Type: CONSTRAINT; Schema: public; Owner: {{DB_USERNAME}}
        --

        ALTER TABLE ONLY public.event
            ADD CONSTRAINT event_pkey PRIMARY KEY (id);
    END IF;

    SELECT sequence_name INTO found_sequence FROM information_schema.sequences WHERE sequence_name = 'hibernate_sequence';
    IF NOT FOUND THEN
        --
        -- Name: hibernate_sequence; Type: SEQUENCE; Schema: public; Owner: {{DB_USERNAME}}
        --

        CREATE SEQUENCE  IF NOT EXISTS public.hibernate_sequence
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;

        EXECUTE 'ALTER TABLE public.hibernate_sequence OWNER TO ' || user_name;
        --ALTER TABLE public.hibernate_sequence OWNER TO user_name;
    END IF;
    
    RETURN 1;

END;
$$;

SELECT db_update('{{DB_USERNAME}}');
