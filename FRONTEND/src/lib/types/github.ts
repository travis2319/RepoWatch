export interface Repo {
  name: string;
  private: boolean;
  language?: string;
}

export interface AccessRepo {
  repo: string;
  hasAccess: boolean;
  language?: string;
}

export interface GithubUser {
  login: string;
  name: string;
  avatar_url: string;
  html_url: string;
  id: number;
}

export interface RateLimit {
  limit: string;
  remaining: string;
  used: string;
  reset: string;
  resource: string;
}

export interface ValidateData {
  valid: boolean;
  user: GithubUser;
  rate_limit: RateLimit;
  token_expiry: string;
}