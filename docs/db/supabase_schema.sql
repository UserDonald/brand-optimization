-- Enable RLS on all tables for tenant isolation
alter table public.organizations enable row level security;
alter table public.competitors enable row level security;
alter table public.competitor_metrics enable row level security;
alter table public.personal_metrics enable row level security;
alter table public.audience_segments enable row level security;
alter table public.content_formats enable row level security;
alter table public.scheduled_posts enable row level security;

-- Organizations (Tenants)
create table public.organizations (
  id uuid not null primary key,
  name text not null,
  account_owner text,
  tier text not null default 'standard',
  created_at timestamptz not null default now()
);

-- Create policy for organizations
create policy "Users can only view their own organization"
on public.organizations
for select
using (auth.uid() = id);

-- Competitors
create table public.competitors (
  id uuid not null default uuid_generate_v4() primary key,
  tenant_id uuid not null references public.organizations(id),
  name text not null,
  platform text not null,
  created_at timestamptz not null default now()
);

-- Create policy for competitors
create policy "Users can only access their tenant's competitors"
on public.competitors
for all
using (auth.uid() = tenant_id);

-- Competitor Metrics
create table public.competitor_metrics (
  id uuid not null default uuid_generate_v4() primary key,
  tenant_id uuid not null references public.organizations(id),
  competitor_id uuid not null references public.competitors(id),
  post_id text not null,
  likes int not null default 0,
  shares int not null default 0,
  comments int not null default 0,
  click_through_rate float not null default 0,
  avg_watch_time float not null default 0,
  engagement_rate float not null default 0,
  posted_at timestamptz not null,
  created_at timestamptz not null default now(),
  unique(tenant_id, competitor_id, post_id)
);

-- Create policy for competitor metrics
create policy "Users can only access their tenant's competitor metrics"
on public.competitor_metrics
for all
using (auth.uid() = tenant_id);

-- Personal Metrics
create table public.personal_metrics (
  id uuid not null default uuid_generate_v4() primary key,
  tenant_id uuid not null references public.organizations(id),
  post_id text not null,
  likes int not null default 0,
  shares int not null default 0,
  comments int not null default 0,
  click_through_rate float not null default 0,
  avg_watch_time float not null default 0,
  engagement_rate float not null default 0,
  posted_at timestamptz not null,
  created_at timestamptz not null default now(),
  unique(tenant_id, post_id)
);

-- Create policy for personal metrics
create policy "Users can only access their tenant's personal metrics"
on public.personal_metrics
for all
using (auth.uid() = tenant_id);

-- Audience Segments
create table public.audience_segments (
  id uuid not null default uuid_generate_v4() primary key,
  tenant_id uuid not null references public.organizations(id),
  name text not null,
  description text,
  segment_size int,
  engagement_style text,
  created_at timestamptz not null default now()
);

-- Create policy for audience segments
create policy "Users can only access their tenant's audience segments"
on public.audience_segments
for all
using (auth.uid() = tenant_id);

-- Content Formats
create table public.content_formats (
  id uuid not null default uuid_generate_v4() primary key,
  tenant_id uuid not null references public.organizations(id),
  name text not null,
  description text,
  avg_engagement_rate float,
  created_at timestamptz not null default now()
);

-- Create policy for content formats
create policy "Users can only access their tenant's content formats"
on public.content_formats
for all
using (auth.uid() = tenant_id);

-- Scheduled Posts
create table public.scheduled_posts (
  id uuid not null default uuid_generate_v4() primary key,
  tenant_id uuid not null references public.organizations(id),
  content text not null,
  scheduled_time timestamptz not null,
  platform text not null,
  format text not null,
  status text not null default 'pending',
  created_at timestamptz not null default now()
);

-- Create policy for scheduled posts
create policy "Users can only access their tenant's scheduled posts"
on public.scheduled_posts
for all
using (auth.uid() = tenant_id);

-- Create functions for analytics
create or replace function get_competitor_comparison(
  p_tenant_id uuid,
  p_competitor_id uuid,
  p_start_date timestamptz,
  p_end_date timestamptz
) returns table (
  metric_type text,
  competitor_value numeric,
  personal_value numeric,
  ratio numeric
) language sql as $$
  with competitor_metrics as (
    select
      sum(likes) as total_likes,
      sum(shares) as total_shares,
      sum(comments) as total_comments,
      avg(engagement_rate) as avg_engagement_rate,
      avg(avg_watch_time) as avg_watch_time
    from public.competitor_metrics
    where tenant_id = p_tenant_id
      and competitor_id = p_competitor_id
      and posted_at between p_start_date and p_end_date
  ),
  personal_metrics as (
    select
      sum(likes) as total_likes,
      sum(shares) as total_shares,
      sum(comments) as total_comments,
      avg(engagement_rate) as avg_engagement_rate,
      avg(avg_watch_time) as avg_watch_time
    from public.personal_metrics
    where tenant_id = p_tenant_id
      and posted_at between p_start_date and p_end_date
  )
  select 'likes' as metric_type, c.total_likes as competitor_value, p.total_likes as personal_value,
         case when p.total_likes = 0 then 0 else c.total_likes::numeric / p.total_likes end as ratio
  from competitor_metrics c, personal_metrics p
  union all
  select 'shares' as metric_type, c.total_shares as competitor_value, p.total_shares as personal_value,
         case when p.total_shares = 0 then 0 else c.total_shares::numeric / p.total_shares end as ratio
  from competitor_metrics c, personal_metrics p
  union all
  select 'comments' as metric_type, c.total_comments as competitor_value, p.total_comments as personal_value,
         case when p.total_comments = 0 then 0 else c.total_comments::numeric / p.total_comments end as ratio
  from competitor_metrics c, personal_metrics p
  union all
  select 'engagement_rate' as metric_type, c.avg_engagement_rate as competitor_value, p.avg_engagement_rate as personal_value,
         case when p.avg_engagement_rate = 0 then 0 else c.avg_engagement_rate::numeric / p.avg_engagement_rate end as ratio
  from competitor_metrics c, personal_metrics p
  union all
  select 'watch_time' as metric_type, c.avg_watch_time as competitor_value, p.avg_watch_time as personal_value,
         case when p.avg_watch_time = 0 then 0 else c.avg_watch_time::numeric / p.avg_watch_time end as ratio
  from competitor_metrics c, personal_metrics p;
$$; 