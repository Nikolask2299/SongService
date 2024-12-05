CREATE OR REPLACE FUNCTION make_tsvector(text TEXT)
    RETURNS tsvector AS $$
BEGIN 
    RETURN to_tsvector('english', text);
END;
$$ LANGUAGE 'plpgsql' IMMUTABLE;


CREATE INDEX IF NOT EXISTS ind_serch_text ON songs
    USING gin(make_tsvector(text));

