export interface MenuItem {
  id: string;
  name: string;
  parent_id: string | null;
  depth: number;
  order_index: number;
  created_at: string;
  updated_at: string;
  children: MenuItem[];
}

export interface CreateMenuRequest {
  name: string;
  parent_id: string | null;
}

export interface UpdateMenuRequest {
  name: string;
}

export interface MoveMenuRequest {
  new_parent_id: string | null;
  order_index: number;
}

export interface ReorderMenuRequest {
  order_index: number;
}

export interface ApiError {
  error: string;
  message: string;
}
