CREATE TABLE IF NOT EXISTS templates (
    template_id UUID NOT NULL,
    version     INTEGER NOT NULL DEFAULT 1,
    name        VARCHAR(255) NOT NULL,
    channel     VARCHAR(20) NOT NULL,
    subject     VARCHAR(500) NOT NULL DEFAULT '',
    body        TEXT NOT NULL,
    variables   JSONB NOT NULL DEFAULT '[]',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (template_id, version)
);

CREATE INDEX idx_templates_active ON templates (template_id, is_active) WHERE is_active = true;
CREATE INDEX idx_templates_name ON templates (name);
CREATE INDEX idx_templates_channel ON templates (channel);
