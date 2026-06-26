import type { DataSource } from './platform';

export type ProductCategory = 'ship' | 'hotel' | 'tour' | 'play' | 'car' | 'service';
export type OrderStatus = 'pending_use' | 'pending_pay' | 'reserved' | 'completed' | 'refunded' | 'refunding';
export type TicketStatus = 'available' | 'unavailable' | 'refunding';
export type PhaseTwoScenarioId = 'phase2-success' | 'phase2-empty' | 'phase2-unauthorized' | 'phase2-failure';

export interface ProductEntry {
  id: string;
  title: string;
  subtitle: string;
  category: ProductCategory;
  priceLabel: string;
  tag: string;
  imageUrl?: string;
  sourceNote?: string;
  actionText: string;
  route?: string;
}

export interface DestinationEntry {
  id: string;
  title: string;
  subtitle: string;
  category: ProductCategory;
  route?: string;
}

export interface UserSummary {
  displayName: string;
  maskedMobile: string;
  realnameStatus: 'verified' | 'placeholder';
  nextTripLabel?: string;
  availableTicketCount: number;
  recentOrderCount: number;
  isLoggedIn: boolean;
}

export interface OrderSummary {
  id: string;
  title: string;
  orderNo: string;
  status: OrderStatus;
  source: 'generic' | 'island_cruise' | 'ship_ticket' | 'hotel' | 'play';
  travelDateLabel: string;
  quantityLabel: string;
  amountLabel: string;
  hint: string;
  ticketId?: string;
}

export interface TicketSummary {
  id: string;
  status: TicketStatus;
  title: string;
  productTitle: string;
  validDateLabel: string;
  validTimeLabel: string;
  quantityLabel: string;
  maskedCode: string;
  verifyLocation: string;
  notice: string[];
  orderNo: string;
  amountLabel: string;
  source: 'generic' | 'island_cruise' | 'ship_ticket' | 'play';
}

export interface MainShellViewModel {
  dataSource: DataSource;
  scenarioId: PhaseTwoScenarioId;
  user: UserSummary | null;
  destinations: DestinationEntry[];
  products: ProductEntry[];
  orders: OrderSummary[];
  tickets: TicketSummary[];
  errorMessage?: string;
}

export interface BackendOrderRecord {
  id: string;
  product_name: string;
  order_no: string;
  status: OrderStatus;
  source_type: OrderSummary['source'];
  travel_date: string;
  quantity_text: string;
  paid_amount_text: string;
  usage_hint: string;
  ticket_id?: string;
}

export interface BackendTicketRecord {
  id: string;
  status: TicketStatus;
  title: string;
  product_name: string;
  valid_date: string;
  valid_time: string;
  quantity_text: string;
  ticket_code: string;
  verify_location: string;
  notice: string[];
  order_no: string;
  paid_amount_text: string;
  source_type: TicketSummary['source'];
}
