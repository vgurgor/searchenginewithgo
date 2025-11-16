-- Seed 24 contents (12 video, 12 text) with varying published_at
INSERT INTO contents (provider_id, provider_content_id, title, content_type, description, url, thumbnail_url, published_at)
VALUES
('yt','v001','Go Concurrency Patterns','video','Talk about concurrency in Go','https://youtu.be/v001','https://img/1.jpg', NOW() - INTERVAL '1 day'),
('yt','v002','Building REST APIs with Gin','video','Tutorial on Gin framework','https://youtu.be/v002','https://img/2.jpg', NOW() - INTERVAL '2 days'),
('yt','v003','PostgreSQL Indexing','video',NULL,NULL,NULL, NOW() - INTERVAL '3 days'),
('yt','v004','Redis Tips','video','Redis best practices',NULL,NULL, NOW() - INTERVAL '4 days'),
('yt','v005','Docker Multi-stage Builds','video',NULL,NULL,NULL, NOW() - INTERVAL '5 days'),
('yt','v006','Kubernetes Basics','video','Intro to K8s',NULL,NULL, NOW() - INTERVAL '6 days'),
('yt','v007','CQRS & Event Sourcing','video',NULL,NULL,NULL, NOW() - INTERVAL '7 days'),
('yt','v008','gRPC with Go','video','gRPC fundamentals',NULL,NULL, NOW() - INTERVAL '8 days'),
('yt','v009','Caching Strategies','video','Cache invalidation',NULL,NULL, NOW() - INTERVAL '9 days'),
('yt','v010','Observability 101','video',NULL,NULL,NULL, NOW() - INTERVAL '10 days'),
('yt','v011','Clean Architecture','video','Structure projects well',NULL,NULL, NOW() - INTERVAL '11 days'),
('yt','v012','Testing in Go','video','Unit testing basics',NULL,NULL, NOW() - INTERVAL '12 days'),
('md','t001','Advanced SQL Joins','text','Deep dive into joins','https://blog/joins','https://img/t1.jpg', NOW() - INTERVAL '1 day'),
('md','t002','Pagination Strategies','text','Offset vs keyset','https://blog/pagination','https://img/t2.jpg', NOW() - INTERVAL '2 days'),
('md','t003','HTTP Caching','text','ETag and Cache-Control','https://blog/http-cache','https://img/t3.jpg', NOW() - INTERVAL '3 days'),
('md','t004','Rate Limiting','text','Token bucket explained','https://blog/rate-limit','https://img/t4.jpg', NOW() - INTERVAL '4 days'),
('md','t005','Idempotency','text','Designing idempotent APIs','https://blog/idempotency','https://img/t5.jpg', NOW() - INTERVAL '5 days'),
('md','t006','Tracing Basics','text','Distributed tracing','https://blog/tracing','https://img/t6.jpg', NOW() - INTERVAL '6 days'),
('md','t007','Security Headers','text','HSTS, CSP, etc.','https://blog/headers','https://img/t7.jpg', NOW() - INTERVAL '7 days'),
('md','t008','SQL Performance','text','Avoid common pitfalls','https://blog/sql-perf','https://img/t8.jpg', NOW() - INTERVAL '8 days'),
('md','t009','Message Queues','text','Use cases and patterns','https://blog/mq','https://img/t9.jpg', NOW() - INTERVAL '9 days'),
('md','t010','Retries & Backoff','text','Exponential strategies','https://blog/retries','https://img/t10.jpg', NOW() - INTERVAL '10 days'),
('md','t011','gRPC vs REST','text','Trade-offs explained','https://blog/grpc-rest','https://img/t11.jpg', NOW() - INTERVAL '11 days'),
('md','t012','Event-driven Design','text','When to use events','https://blog/events','https://img/t12.jpg', NOW() - INTERVAL '12 days');

-- Seed metrics using provider keys to resolve content ids
INSERT INTO content_metrics (content_id, views, likes, reading_time, reactions, final_score, recalculated_at)
SELECT c.id,
       CASE WHEN c.content_type='video' THEN (random()*10000)::bigint ELSE 0 END AS views,
       CASE WHEN c.content_type='video' THEN (random()*2000)::bigint ELSE 0 END AS likes,
       CASE WHEN c.content_type='text' THEN (10 + (random()*20)::int) ELSE 0 END AS reading_time,
       CASE WHEN c.content_type='text' THEN (random()*500)::int ELSE 0 END AS reactions,
       round((random()*100)::numeric, 2) AS final_score,
       NOW() - INTERVAL '1 hour'
FROM contents c
WHERE NOT EXISTS (SELECT 1 FROM content_metrics cm WHERE cm.content_id = c.id);


