import type { DataSource } from './platform';

export type IslandCruiseStep = 'detail' | 'traveler' | 'pay' | 'ticket';
export type IslandCruiseScenarioId =
  | 'phase3-success'
  | 'phase3-no-voyage'
  | 'phase3-failure'
  | 'phase3-insufficient-stock'
  | 'phase3-lock-expired'
  | 'phase3-passenger-invalid'
  | 'phase3-sale-failed';

export interface IslandPort {
  id: number;
  name: string;
}

export interface IslandFare {
  id: string;
  name: string;
  cabinClassId: string;
  cabinClassName: string;
  cabinTypeCode: string;
  price: number;
  originalPrice: number;
  stock: number;
}

export interface IslandVoyage {
  id: number;
  name: string;
  voyageNo: string;
  shipName: string;
  departureDate: string;
  departureTime: string;
  upPort: IslandPort;
  downPort: IslandPort;
  stateLabel: string;
  fares: IslandFare[];
}

export interface IslandPassenger {
  name: string;
  mobile: string;
  certTypeId: number;
  certTypeName: string;
  certNo: string;
  fareId: string;
}

export interface IslandOrderDraft {
  contactName: string;
  contactMobile: string;
  passengers: IslandPassenger[];
  quantity: number;
  totalAmount: number;
}

export interface IslandLockedOrder {
  localOrderNo: string;
  supplierOrderNo: string;
  ticketNo: string;
  expireTimeLabel: string;
  status: 'pending_payment' | 'expired';
  amount: number;
}

export interface IslandTicketResult {
  localOrderNo: string;
  ticketNo: string;
  codeContent: string;
  maskedCode: string;
  paidAtLabel: string;
  voyageLabel: string;
  passengerLabel: string;
  status: 'ticketed' | 'sale_failed';
  failureMessage?: string;
}

export interface IslandCruiseViewModel {
  dataSource: DataSource;
  scenarioId: IslandCruiseScenarioId;
  hero: {
    title: string;
    summary: string;
    imageUrl: string;
  };
  ports: {
    up: IslandPort;
    down: IslandPort;
  };
  certTypes: Array<{ id: number; name: string }>;
  recommended: {
    date: string;
    label: string;
    firstTime: string;
    count: number;
    minPriceLabel: string;
  } | null;
  voyages: IslandVoyage[];
  selectedVoyageId: number | null;
  draft: IslandOrderDraft;
  lockedOrder: IslandLockedOrder | null;
  ticket: IslandTicketResult | null;
  serviceNotes: string[];
  errorMessage?: string;
}

export interface BackendIslandFareRecord {
  fareTypeId?: number | string;
  fareTypeName?: string;
  price?: number | string;
  settlePrice?: number | string;
  originalPrice?: number | string;
  realName?: number;
}

export interface BackendIslandCabinRecord {
  cabinClassId?: number | string;
  cabinClassName?: string;
  cabinTypeCode?: string;
  cabinType?: string;
  count?: number | string;
  fareTypeList?: BackendIslandFareRecord[];
}

export interface BackendIslandVoyageRecord {
  voyageId: number;
  voyageName?: string;
  lineName?: string;
  voyageNo?: string;
  shipName?: string;
  departureDate?: string;
  departureTime?: string;
  upTime?: string;
  upPortId?: number;
  upPortName?: string;
  downPortId?: number;
  downPortName?: string;
  voyageState?: string;
  cabinPriceList?: BackendIslandCabinRecord[];
}

export interface BackendIslandLockRecord {
  local_order_no?: string;
  orderNo?: string;
  ticketNo?: string;
  lock_expire_at?: string;
  amount?: number;
}

export interface BackendIslandSaleRecord {
  local_order_no?: string;
  ticketNo?: string;
  codeContent?: string;
  paid_at?: string;
}
