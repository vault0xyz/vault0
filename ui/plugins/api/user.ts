import type {
  ICreateUserRequest,
  IPagedResponse,
  IUpdateUserRequest,
  IUser,
} from '~/types';
import { User } from '~/types';
import type { ApiClient } from './client';
import { API_ENDPOINTS } from './endpoints';

/**
 * Client for interacting with user-related API endpoints
 */
export class UserClient {
  private client: ApiClient;

  /**
   * Creates a new user client
   * @param client API client instance
   */
  constructor(client: ApiClient) {
    this.client = client;
  }

  /**
   * Creates a new user
   * @param request User creation request
   * @returns Created user
   */
  async createUser(request: ICreateUserRequest): Promise<IUser> {
    const data = await this.client.post<any>(API_ENDPOINTS.USERS.BASE, request);
    return User.fromJson(data);
  }

  /**
   * Updates a user's properties
   * @param id User ID
   * @param request User update request
   * @returns Updated user
   */
  async updateUser(id: string, request: IUpdateUserRequest): Promise<IUser> {
    const endpoint = API_ENDPOINTS.USERS.BY_ID(id);
    const data = await this.client.put<any>(endpoint, request);
    return User.fromJson(data);
  }

  /**
   * Deletes a user
   * @param id User ID
   */
  async deleteUser(id: string): Promise<void> {
    const endpoint = API_ENDPOINTS.USERS.BY_ID(id);
    await this.client.delete(endpoint);
  }

  /**
   * Gets a user by ID
   * @param id User ID
   * @returns User details
   */
  async getUser(id: string): Promise<IUser> {
    const endpoint = API_ENDPOINTS.USERS.BY_ID(id);
    const data = await this.client.get<any>(endpoint);
    return User.fromJson(data);
  }

  /**
   * Lists users with token-based pagination
   * @param limit Maximum number of users to return (default: 10)
   * @param nextToken Token for retrieving the next page of results (default: undefined)
   * @returns Paginated list of users
   */
  async listUsers(limit: number = 10, nextToken?: string): Promise<IPagedResponse<IUser>> {
    const params: Record<string, any> = { limit };
    if (nextToken) {
      params.next_token = nextToken;
    }
    
    const data = await this.client.get<any>(API_ENDPOINTS.USERS.BASE, params);
    return {
      items: User.fromJsonArray(data.items || []),
      limit: data.limit,
      nextToken: data.nextToken
    };
  }
} 