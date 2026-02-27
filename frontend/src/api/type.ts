export interface NekosImageData {
  id: string;
  category: string;
  colors: {
    main: string;
    palette: string[];
  };
  image: {
    original: {
      url: string;
      extension: string;
    };
    compressed: {
      url: string;
      extension: string;
    };
  };
  metadata: {
    original: {
      width: number;
      height: number;
      size: number;
      extension: string;
    };
    compressed: {
      width: number;
      height: number;
      size: number;
      extension: string;
    };
  };
  rating: "safe" | "suggestive";
  tags: string[];
  anime: {
    title: string | null;
    character: string | null;
  };
  source: {
    url: string | null;
    direct: string | null;
  };
  attribution: {
    artist: {
      username: string | null;
      profile: string | null;
    };
    copyright: string | null;
  };
}

export interface ResponseGetGirls {
  success: boolean;
  status: number;
  count: number;
  images: NekosImageData[];
}
