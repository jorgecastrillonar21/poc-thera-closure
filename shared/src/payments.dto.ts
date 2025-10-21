export interface CreateCheckoutSessionDTO {
  priceId: string;
  successUrl: string;
  cancelUrl: string;
  userId: string;
}

export interface CheckoutSessionResponseDTO {
  sessionId: string;
  url: string;
}

export interface StripeEventDTO {
  id: string;
  type: string;
  data: {
    object: any;
  };
}

export interface SubscriptionDTO {
  id: string;
  customerId: string;
  status: string;
  currentPeriodStart: Date;
  currentPeriodEnd: Date;
  priceId: string;
  productId: string;
}

export interface CustomerPortalDTO {
  userId: string;
  returnUrl: string;
}

export interface CustomerPortalResponseDTO {
  url: string;
}