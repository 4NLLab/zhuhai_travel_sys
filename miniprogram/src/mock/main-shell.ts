import type {
  BackendOrderRecord,
  BackendTicketRecord,
  DestinationEntry,
  MainShellViewModel,
  PhaseTwoScenarioId,
  ProductEntry,
  UserSummary
} from '@/types/main-shell';
import { mapBackendOrder, mapBackendTicket } from '@/utils/main-shell-mappers';

const sharedDestinations: DestinationEntry[] = [
  { id: 'island-cruise', title: '澳门环岛游', subtitle: '海上观光', category: 'ship', route: '/pages/island-cruise/index' },
  { id: 'shekou-ferry', title: '深圳香港船票', subtitle: '九洲蛇口', category: 'ship' },
  { id: 'private-tour', title: '私人旅游订制', subtitle: '专属路线', category: 'tour' },
  { id: 'bridge', title: '港珠澳大桥游', subtitle: '湾区观光', category: 'tour' },
  { id: 'chimelong', title: '珠海长隆', subtitle: '亲子度假', category: 'play' },
  { id: 'macau', title: '澳门游', subtitle: '轻松过关', category: 'tour' },
  { id: 'hotel', title: '珠海酒店', subtitle: '湾区住宿', category: 'hotel' },
  { id: 'car', title: '口岸接送', subtitle: '专车出行', category: 'car' }
];

const sharedProducts: ProductEntry[] = [
  {
    id: 'island-cruise-night',
    title: '九洲港至蛇口船票',
    subtitle: '对接九洲港分销接口，真实航班、票价和锁票链路',
    category: 'ship',
    priceLabel: '¥140 起',
    tag: '热卖',
    imageUrl: '/static/phase2/macau-cruise-night-banner-web.jpg',
    actionText: '预订',
    route: '/pages/island-cruise/index',
    sourceNote: '来源参考 index.html 首页船票推荐'
  },
  {
    id: 'hotel-package',
    title: '珠海横琴酒店双人套餐',
    subtitle: '近长隆与口岸，含早餐，可加购接送服务',
    category: 'hotel',
    priceLabel: '¥399 起',
    tag: '套餐',
    imageUrl: '/static/phase2/zhuhai-bay-home-hero-web.jpg',
    actionText: '预订',
    sourceNote: '非首版交易路径占位'
  },
  {
    id: 'hk-macau-route',
    title: '香港澳门一日游路线',
    subtitle: '口岸集合，可选导游、用车、门票与酒店组合',
    category: 'tour',
    priceLabel: '¥268 起',
    tag: '路线',
    imageUrl: '/static/phase2/taxi-scan-illustration-web.jpg',
    actionText: '咨询',
    sourceNote: '非首版交易路径占位'
  },
  {
    id: 'chimelong-ticket',
    title: '长隆门票与石景山缆车',
    subtitle: '亲子玩乐组合，支持单品或套票售卖',
    category: 'play',
    priceLabel: '¥128 起',
    tag: '电子票',
    imageUrl: '/static/phase2/ticket-wallet-illustration-web.jpg',
    actionText: '预订',
    sourceNote: '非首版交易路径占位'
  },
  {
    id: 'port-transfer',
    title: '珠港澳口岸接送租车',
    subtitle: '商务车、包车半日游、机场码头接送',
    category: 'car',
    priceLabel: '¥180 起',
    tag: '接送',
    imageUrl: '/static/phase2/taxi-scan-illustration-web.jpg',
    actionText: '询价',
    sourceNote: '非首版交易路径占位'
  }
];

const userSummary: UserSummary = {
  displayName: '珠海游客',
  maskedMobile: '138****6321',
  realnameStatus: 'verified',
  nextTripLabel: '下一段行程：澳门环岛游夜景船票 · 2026-07-15 19:30',
  availableTicketCount: 3,
  recentOrderCount: 5,
  isLoggedIn: true
};

const backendOrders: BackendOrderRecord[] = [
  {
    id: 'order-island-1',
    product_name: '澳门环岛游夜景船票',
    order_no: 'ZH20260715001028',
    status: 'pending_use',
    source_type: 'island_cruise',
    travel_date: '2026-07-15',
    quantity_text: '成人票 x1',
    paid_amount_text: '¥88',
    usage_hint: '凭电子票码到湾仔旅游码头核销',
    ticket_id: 'ticket-island-1'
  },
  {
    id: 'order-hotel-1',
    product_name: '横琴酒店双人套餐',
    order_no: 'ZH20260624000995',
    status: 'reserved',
    source_type: 'hotel',
    travel_date: '2026-06-24',
    quantity_text: '1 间 1 晚',
    paid_amount_text: '¥399',
    usage_hint: '到店凭证件办理入住'
  },
  {
    id: 'order-car-1',
    product_name: '珠港澳口岸接送租车',
    order_no: 'ZH20260626001001',
    status: 'pending_pay',
    source_type: 'generic',
    travel_date: '2026-06-26',
    quantity_text: '商务车 x1',
    paid_amount_text: '待支付 ¥180',
    usage_hint: '请在保留时间内完成支付'
  },
  {
    id: 'order-tour-1',
    product_name: '香港澳门一日游路线',
    order_no: 'ZH20260618000842',
    status: 'completed',
    source_type: 'generic',
    travel_date: '2026-06-18',
    quantity_text: '2 成人',
    paid_amount_text: '¥536',
    usage_hint: '可查看行程回顾与发票'
  },
  {
    id: 'order-play-1',
    product_name: '长隆门票与石景山缆车',
    order_no: 'ZH20260619000711',
    status: 'refunded',
    source_type: 'play',
    travel_date: '2026-06-19',
    quantity_text: '2 张套票',
    paid_amount_text: '¥256',
    usage_hint: '退款已原路退回'
  }
];

const backendTickets: BackendTicketRecord[] = [
  {
    id: 'ticket-island-1',
    status: 'available',
    title: '待使用券码',
    product_name: '澳门环岛游夜景船票',
    valid_date: '2026-07-15',
    valid_time: '19:30-21:00',
    quantity_text: '成人票 x1',
    ticket_code: '0105049143806',
    verify_location: '湾仔旅游码头检票口',
    notice: ['核销前请勿将券码提供给无关人员。', '如遇天气或航班调整，请以码头现场通知为准。'],
    order_no: 'ZH20260715001028',
    paid_amount_text: '¥88',
    source_type: 'island_cruise'
  },
  {
    id: 'ticket-bridge-unavailable',
    status: 'unavailable',
    title: '不可用券码',
    product_name: '港珠澳大桥游电子票',
    valid_date: '2026-06-20',
    valid_time: '10:00-12:00',
    quantity_text: '成人票 x1',
    ticket_code: 'BRIDGE20260620001',
    verify_location: '九洲港服务台',
    notice: ['该票券已过有效期或被售后锁定。', '如需继续使用，请联系平台客服处理。'],
    order_no: 'ZH20260620000722',
    paid_amount_text: '¥128',
    source_type: 'ship_ticket'
  }
];

function successScenario(): MainShellViewModel {
  return {
    dataSource: 'mock',
    scenarioId: 'phase2-success',
    user: userSummary,
    destinations: sharedDestinations,
    products: sharedProducts,
    orders: backendOrders.map(mapBackendOrder),
    tickets: backendTickets.map(mapBackendTicket)
  };
}

export function normalizePhaseTwoScenario(scenarioId: string): PhaseTwoScenarioId {
  if (
    scenarioId === 'phase2-empty' ||
    scenarioId === 'phase2-unauthorized' ||
    scenarioId === 'phase2-failure' ||
    scenarioId === 'phase2-success'
  ) {
    return scenarioId;
  }
  return 'phase2-success';
}

export function getMainShellMock(scenarioId: string): MainShellViewModel {
  const normalizedScenarioId = normalizePhaseTwoScenario(scenarioId);
  if (normalizedScenarioId === 'phase2-empty') {
    return {
      ...successScenario(),
      scenarioId: normalizedScenarioId,
      orders: [],
      tickets: [],
      user: {
        ...userSummary,
        availableTicketCount: 0,
        recentOrderCount: 0,
        nextTripLabel: undefined
      }
    };
  }

  if (normalizedScenarioId === 'phase2-unauthorized') {
    return {
      ...successScenario(),
      scenarioId: normalizedScenarioId,
      user: null,
      orders: [],
      tickets: [],
      errorMessage: '请先登录后查看订单和票券。'
    };
  }

  if (normalizedScenarioId === 'phase2-failure') {
    return {
      ...successScenario(),
      scenarioId: normalizedScenarioId,
      orders: [],
      tickets: [],
      errorMessage: '订单服务暂时不可用，请稍后重试。'
    };
  }

  return successScenario();
}
