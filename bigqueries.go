package main

const WeeklyUpdateQuery = `
SELECT
  type, JSON_EXTRACT(payload, '$.action') as action, JSON_EXTRACT(payload, '$.merged') as merged
FROM
  TABLE_DATE_RANGE(githubarchive:day.events_, TIMESTAMP('2015-03-01'), TIMESTAMP('2015-03-27'))
WHERE
  actor.login ='dhh' AND type = 'PullRequestEvent';
`

const HistoricArchiveQuery = `
SELECT
  concat( string(payload_pull_request_id), "-", payload_pull_request_user_login, "-", payload_pull_request_merged_by_login) as id,
  payload_action,
  payload_pull_request_merged,
  payload_pull_request_title,
  payload_pull_request_url,
  payload_pull_request_user_login,
  payload_pull_request_merged_by_login,
  payload_pull_request_merged_at
FROM
  TABLE_QUERY(githubarchive:month, 'YEAR(CONCAT(table_id, "01")) < 2015')
WHERE
  payload_action = 'closed' AND
  payload_pull_request_merged = 'true' AND
  payload_pull_request_user_login IN (%s) OR
  payload_pull_request_merged_by_login IN (%s)
LIMIT 1000
`
