CREATE TABLE users (
   id         text PRIMARY KEY,
   name       text NOT NULL,
   is_active  boolean NOT NULL DEFAULT true,
   team_name  text
);

CREATE TABLE teams (
   name text PRIMARY KEY
);

ALTER TABLE users
    ADD CONSTRAINT users_team_fk
    FOREIGN KEY (team_name) REFERENCES teams(name);

CREATE TABLE pull_requests (
   id         text PRIMARY KEY,
   name       text NOT NULL,
   author_id  text NOT NULL,
   status     text NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
   created_at timestamptz NOT NULL DEFAULT now(),
   merged_at  timestamptz
);

ALTER TABLE pull_requests
    ADD CONSTRAINT pull_requests_author_fk
    FOREIGN KEY (author_id) REFERENCES users(id);

CREATE TABLE pull_request_reviewers (
    pull_request_id text NOT NULL,
    reviewer_id     text NOT NULL,
    PRIMARY KEY (pull_request_id, reviewer_id)
);

ALTER TABLE pull_request_reviewers
    ADD CONSTRAINT pr_reviewers_pr_fk
    FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE;

ALTER TABLE pull_request_reviewers
    ADD CONSTRAINT pr_reviewers_user_fk
    FOREIGN KEY (reviewer_id) REFERENCES users(id);