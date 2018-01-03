package shortuuid

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

var testVector = []struct {
	uuid      string
	shortuuid string
}{
	{"c3eeb3e6-e577-4de2-b5bb-08371196b453", "nSKHvGZM4M5A4KyN4zp4sc"},
	{"7c37b0a5-9d63-4d91-a79e-85a480fd2348", "gos3umowEP9bS8QAGSuk7Q"},
	{"52600a3d-a12a-4041-aa9c-afd5e5923a70", "ZHhWmgrdXxgMqzgYRXGRfG"},
	{"7b80853c-a9f3-4acd-b735-17d29b3b9122", "gGpmQT6oVMVeQvnLkgJWyP"},
	{"2f56ebd5-373f-478f-afab-1142b74b75bc", "Hi7d4MYcbjKxzeg8YiN7SA"},
	{"72573478-ade8-4a6c-a583-d098e9c24de3", "yoANdrf88xUXvwbS5GRbMN"},
	{"f4c03714-aa7e-42e3-91f8-9a7a3bc05f75", "VTtAQPXV7MRyAghQPA6AZm"},
	{"f9ee01c3-2015-4716-930e-4d5449810833", "7fm8A5kT3j9j5M2HcjofUn"},
	{"a5fb9251-9373-486b-b9bd-ecf1c7cc6d9c", "DQqWKCLAjsm6fTTFCwuKYX"},
	{"9cbfe3da-dc50-4d07-9d3a-3d161e5a82b1", "dhRv9nyctV3M6uUJtcZgtV"},
	{"cda34e96-188b-41b2-acac-05cd288c32e1", "Cekw67uyMpBGZLRP2HFVbe"},
	{"0ce70d9e-830e-404b-8b56-25cc611eed6a", "qXPoF8kNjxoqWZVCnoXrJ4"},
	{"36a06807-8c07-4e65-8b99-855b3628c155", "LXuzZwRWpEYos2yAiDmziB"},
	{"4672c1dc-125a-4469-a81e-c69dc3b8f792", "p6tjrWt5UH9RpAUqPfjTYE"},
	{"e63d7da2-7f8d-4954-9c2d-d46b324d10f2", "P3bYrwnrsvRXdsRCTByyxi"},
	{"926251a7-41a3-4029-bf44-bf860e881423", "r3FdNDkxMf4rdz8V7dcZ4U"},
	{"464a919f-54e9-435e-a029-a3b942b07134", "THvwuXRQt9soLzL3L2zrWE"},
	{"7c24ef71-f4eb-4000-afb6-710c7f9e0942", "9tXMVg3xaJxwuK7v4XZ27Q"},
	{"d3973d0f-22e4-44b5-a03f-34a57458a8ff", "rqcFhGVXy7quTXiERRPref"},
	{"74102102-3d04-48fc-b2b5-6ef50109e27f", "v8tx4KDUXyfSFcSmFZ45fN"},
	{"67cdb44b-f0fb-47c1-9060-025216c14067", "nD4n6pXuc9BanbiQM6HjUL"},
	{"503c91b0-6a52-4b26-be71-b96c1d4d908b", "fq2q5eoUmgGuuKTtCa3jHG"},
	{"174ea6fa-9321-4e80-b732-006b11998663", "zw3YBp3erLmWGWFvNW8PA6"},
	{"fa60d3a9-df24-468c-8006-e2a4f37ddf17", "zohky26otniDxSonDD6EZn"},
	{"eff670a0-75f5-4458-afd4-9aa442dcc507", "feViSxBJPp7ZUgCA2hCbhk"},
	{"2c301095-3f99-4fc8-88e3-a2be83fd39fc", "YAJdktpVp9ioZXyerGS9s9"},
	{"ee2d3de2-4ad7-4f32-9071-f8aa4b46dc0b", "XrKarC8cp4dWSRnWQjoUPk"},
	{"12942a1b-cf39-44d6-bfbb-0ae17affd408", "tYx4TiGiu78HcFfkmgkRK5"},
	{"f9b2ea09-32c6-426b-bb3b-a269182e99f3", "b7rFanuSDdkVW4Af7vNLSn"},
	{"02f3da16-b995-445a-877b-2319f0c06b41", "5aQC5t9BRRcicH9GbfmwX"},
	{"b8538991-0957-4bab-93f3-7f78ca9c1550", "vDJmAtHGFpqmrP6MNTjMoa"},
	{"3bf87639-ba74-463e-9455-bf7bcb8fd898", "4pUYNRFHTG3YVgThPZvCgC"},
	{"c9c7cafa-a131-4193-a21d-b65d0052f77b", "ziGxvYy3dYbecxB2VPNNud"},
	{"f921ebe3-2d1b-42fc-951b-7b275d5a1e5f", "pwt3bwxM4WTMGRqNTuxaLn"},
	{"c3d5e2c5-1bc8-4a86-96fb-bf0877c9fc03", "UNqNwmt3Eb9TDXRJXon5rc"},
	{"582dd24f-beda-42cf-b751-7c7e9827b7d1", "cjhuJfBhZQiSZ7GkYPGHhH"},
	{"6dfacd19-5cfd-4b95-92be-7e9fb0fe343e", "C6wrxXg46LWXAnhNTyVNaM"},
	{"ea3f3050-d7ce-4ad7-ae24-9b7b164dcec1", "r29MuyWbYNZ75FhWmh5dgj"},
	{"07dcc180-1323-4985-8359-296261dd53c8", "Qw7DZuPhurJbWAAGwhwjQ3"},
	{"89ffb519-4cea-4cd2-b051-0da7b451ca1a", "q7C7JQRar8QxqhkePNRXZS"},
	{"81fa1f63-389e-44ae-906d-4c1f4a715a61", "mQVvJQHRPinBYydtDaHB9R"},
	{"77eee831-5e48-400b-b1e1-31bc8d28876e", "hnGWZSNK7xyhivxQoYJKMP"},
	{"63815548-e46f-441a-8bdf-3153db75e56b", "Qnmrn7WRTv2kmDC5oNZ9iK"},
	{"60cba957-dc15-404a-994a-cab6b11b1315", "po3Y5UByzV7T8ny9fvCfEK"},
	{"a23360e0-1caf-4275-b957-0bcd7aa109d8", "mRV85qM4gaVEVidUCrfxrW"},
	{"48610926-01fb-4ef8-ae6b-ba9b45569525", "4YcLEHE6evYdnN9R6Er3tE"},
	{"c1d07707-c751-4e18-bf53-924a02e7a312", "tfdCPzDgjH7kDZmoUbRaVc"},
	{"d22b3035-2213-49a8-8399-a547303f29a3", "9K49a86KhGxzjfdt8YLSQf"},
	{"ba40a392-f1f6-4486-8c27-00b78db13ddc", "iHHZerT4hpfGqMbjuXCt9b"},
	{"241d6b4e-81a7-439a-896d-b645ec223eb5", "fuTJdCqF7jvgn6p8ZRoGS8"},
	{"e2f1fc6f-c8c8-4e6c-91ff-aca0ca558535", "qHhWwrK3ZPHZ3w5LDjHaPi"},
	{"b6f1e412-da4f-43bf-8a9f-ae4658f5b9ce", "7EohUgpsHtHugJudodBMZa"},
	{"db850c2c-a4e2-45c5-9b67-fda3be319df2", "uAQEWZNR9MQrqTVCU3qG5h"},
	{"63ac6e18-e0c9-477d-b08f-039b9538e577", "Qh3dcsgZBQCseB34xNsqjK"},
	{"eda9eacb-d698-41af-8790-4f26ecbd3a39", "29XrEqVWgteJF2Kr6xGHJk"},
	{"76d61350-f831-40fe-ae88-c2d63885c954", "u6xS5LyehZPedbnBRjACAP"},
	{"9dc996f5-2dfb-4e60-96a2-f11f59eb71de", "B8tRdDRCz6spwRDVPoXD6W"},
	{"9e0ef764-f98d-4c66-abd8-5d5e7392bacb", "TsnXJsfNMZGvuu9WJHCx8W"},
	{"08c79281-f45e-4463-b26d-1dc84ece1a05", "VgbLFBmKXHXs5VZdxuB4a3"},
	{"6611ce43-fe22-47ef-879e-de3fbbd4d9f6", "HHsWgxyECxDUUdNKEvu9BL"},
	{"748374a4-43c4-42a4-92be-f70daffe5d58", "h6793J2TNezTRPo8iFUdjN"},
	{"eb571f9d-dbe4-4d37-840f-0db7c9a35e3a", "snrX2R8WkCAALQ5Ux6Cisj"},
	{"dbc3f7d2-5f4b-4fb9-b473-2ca74cf7e2ef", "WRGi6eUnSKwVAaCwcXuk7h"},
	{"4e4d65f0-22b1-48ba-acc4-b378e7a826e3", "ybEmP4TsGe53BYtqDAu7wF"},
	{"c5232425-c634-4653-a367-0eb07235ccfc", "U6XjtF8PADo2bW7Gd6KH6d"},
	{"8a95bedb-0592-4329-87be-8affdf5f2e97", "ck493bWVpcBsppo94mDUfS"},
	{"52fa1417-02e9-4656-96ec-5cf44fdcd3c2", "RyF2UdmNpgzbxjNVbo6XmG"},
	{"d728fb9b-240b-4ca1-b864-50b807364acb", "y3kctB9VyUVeDABHVQg4Jg"},
	{"9da4fb93-9260-46a7-b6df-f04218560fe0", "jxBVLETqKeFjofyJe8sk4W"},
	{"fa9c1ef1-1f4d-4b74-9590-878c79f1f353", "xKpbsbfBxJA5r5uC7wxZbn"},
	{"64d1355f-d052-4bd9-83f4-39b93fb1c01f", "fCd2nkGhNF9UBcxpDsySwK"},
	{"7227d320-fd77-4430-8e9c-ac3f93e07c51", "crnPpyJpmterz6sLUwRiKN"},
	{"22eac9a7-b0b7-4c8b-9e88-4a2c4a48eef4", "PKMNKhycsviQJbcG3uQ8E8"},
	{"68252941-638b-4e0d-807a-3e8cd04f2493", "iRPfmXvLunFTvtoanekCYL"},
	{"3c9886b3-cabc-4854-8fe3-d394c0901830", "ELPAWnjUbWpbS3xhJSNYnC"},
	{"37eea0c0-ec18-42e8-a152-c98fee7346fb", "QgZm5asemZJYKqpCJrTFxB"},
	{"a61daaac-ea00-4bba-89d9-22b97763a43a", "aCfFqp3vmiADaHuL6LufZX"},
	{"1fc9211b-c1b8-464e-ae01-1c59d1972f47", "ENBaEjcoJktWRqfBkVDNf7"},
	{"21e24536-3b3e-4166-85c4-c29ea151fbbe", "LNZqphXG6tigmWJeBt7e38"},
	{"94fac6e4-315c-4399-89c0-0c5bdd5bf206", "w6qsE9trfLSjwUMjzxzsWU"},
	{"e9ae9ba7-4fb1-4a6d-bbca-5315ed438371", "vACGXSMs5eYcfXgAAjbtaj"},
	{"39439367-6f10-4de8-883f-7ace28c3bc33", "KSbYQEttmBpDjSkUxEMkCC"},
	{"ac6250d9-6e12-4686-98b2-3e5e8ab11ecb", "ZZVXrsGQdVqLz4JGNdKFgY"},
	{"f1975e43-a054-4243-85b1-14ed27be8330", "Bb45oofgN75sPKJYWWe7zk"},
	{"9206e127-8559-4cc4-877e-5cf342aecd25", "7LppjN8wAobe3FXybT9wyT"},
	{"4aeff9a2-4ad2-4c85-abd1-95e798b2132a", "JPu3vSYW3LEGp7RTMSkyLF"},
	{"c7c541b0-d024-432e-8b84-dcc011bfb596", "ag4yKsYNczpg2gKA3TWyYd"},
	{"852fcc13-5632-4188-92a8-12909c77f83f", "tJusodB6eZxNyVrGXGgihR"},
	{"97d87c66-f983-48b7-84c4-1d521e790ecf", "jYFL8gwyq9inicFz4dmw2V"},
	{"4f800ff2-7c1a-4305-b7d5-ce4c5b08e928", "wGSD5xRu6Kc9z8PrFsNGAG"},
	{"185c8229-dbd0-4fc1-a485-04a0532a71d0", "huiUUTVakc3eBxVmQgU5M6"},
	{"0170eb09-dd9d-4682-a58f-b49eb481a40d", "YE6XcuFs8SBnRhgVcS4dG"},
	{"d11c27c7-e302-4007-8a82-a83f57da217a", "uRQsP9P54sbnWjw66uKhDf"},
	{"4ef326fe-7886-4054-b427-4cd9f48200bf", "ibgFYZiZccqvy2PrsMCg4G"},
	{"9e5366f3-6bb8-42db-b999-1dc3612be8e4", "nWdjqVaVdyt94U794eifBW"},
	{"75a8dcb5-9e11-4fbb-895b-60bb6959b6e7", "9tvqGWgZgfiq2mPAvazFwN"},
	{"2b3bd265-a475-401d-88a6-3aa8580483ba", "zdGUGgimEzGkL65LnmuTh9"},
	{"7273ea84-04a1-4d32-92ff-b0afad3eedb6", "6FBBvgSqZeKBpyxxBfFjNN"},
	{"d2af90cd-1b97-4bec-810b-562406b31470", "aTfqZJoFH5kXHiYGxoFgVf"},
	{"8ec76efb-b891-4ee4-868d-997597e390e1", "mdoaq3unkSKHuuDKc6gzQT"},
	{"431933fb-14dd-411c-8fc6-5ca57fa7979d", "ajLWxEodc6CmQLHADuKVwD"},
	{"4788fd99-a1f3-468a-bc67-ada5f80016f4", "QvqUofYurTwpHNWUG52VjE"},
	{"502f56e0-f351-4f42-b395-4dde8b8c19b6", "B8tAKsW9Zw6ZRHySjmADHG"},
	{"b3150123-4e7c-4080-a029-c7437c6513df", "JYqm5dazak2ye8hfs7DBsZ"},
	{"61714a37-a894-4ec5-a569-7e925b46e368", "no6QHS2MfYKUZ6ys4xCEMK"},
	{"8a2f74ed-8146-4f86-a50a-aa28a21c3829", "t3HZ5drXDpPRgjuEnBFRbS"},
	{"5ea5a1a8-bd7e-46fd-ae88-e2d779bc9453", "uQod4vaUSE3ciZdwbbCsqJ"},
	{"35165d02-262d-4e30-8344-c70ee9cc63ee", "DjB4ac8mxbB25FCQswyPTB"},
	{"0026636a-e9b3-4a88-9c66-bf49d8cad81f", "BcRZF3N5Bv7FcJTEsxgX3"},
	{"84973303-f7ff-47a0-8b47-c54767700622", "kzcnTJSM4nPRXJT89Q7gbR"},
	{"92faac68-0387-49ae-a219-45e7c13fd6a3", "Rw8KtoWPAvbt2EqPtAebAU"},
	{"e9f17c46-0d3d-45e8-bc36-6c793b59c91c", "QXh7kFCLxJkJjayBAXcYdj"},
	{"d39295c8-f4d8-407f-bcdf-d07911f4d997", "tWkeanmnjjupCMcjnUsfef"},
	{"a2952f3b-fba2-4926-b19d-26d550093b93", "VEaZCnxWXK6VryEoLUXqvW"},
	{"40e86d7b-77dd-4f55-aec8-29262e028c38", "3Ad8xKbmR43ud53T7Y4HZD"},
	{"c6876847-fda6-4ab5-a2c2-9d8fe9c543bc", "Bj5frkc73iykq7XVT6nPLd"},
	{"cf5c7de8-82d2-475f-a078-825b76d0a287", "PzWDnpB8hSu8cZN38AUxte"},
	{"c4560036-bd08-42d6-ba59-1a49a51643df", "zYWYB5k4sQ55k4KL7V6Awc"},
	{"ded85fd3-afe7-45a1-a4d8-90c0a5bc8901", "MDkC9BQWzG7RPvHpCGCzeh"},
	{"f8221ca8-49e2-44e5-a104-7cec6fb28d48", "Q682Tt3m3ohTTdoGZgLTAn"},
	{"c1544cac-c6ee-40a7-98c8-900a7282ae22", "43nxCsa2kHWeAK32kG4fQc"},
	{"422a2339-0950-4a68-aec1-8c4056feef2a", "SBxqULeq9wXZskuxLoV2nD"},
	{"d9f26801-08d7-42e2-adee-e178944d9ba9", "z2QKKXd38qApcjuav6eKng"},
	{"69807fd5-8bcf-4d4c-a534-ff298d4a5a0a", "fe9aUtsDhyC9ssuU2Y5xmL"},
	{"65656f1e-ae8d-48a3-825d-4499144eefab", "bGSD7rg3Xm2Y4mjNDygK4L"},
	{"0a4399e6-7fad-452b-acc8-1d6b33415019", "NzVFQCNbRxxccPCpHRK7q3"},
	{"4860e725-08f5-4ba1-87e6-1db791f183d8", "8Ka6kDcVRxniasWrb8Z3tE"},
	{"709a8e5a-6320-4649-b4c2-8aceb28502ae", "ekLbSTDNHj9XSCjo6XNy3N"},
	{"3b6aa689-ab23-4ded-932d-090e029a7a99", "4PjpHyzbdreNeUxSc3iaaC"},
	{"00302322-9d08-4394-a1ac-b25bbab774d4", "zv3wa6Hm2fK2SaaRLkhu3"},
	{"12a12ca0-0d6e-4b84-bd6f-9715d467ae76", "iCLCA6LpyQKnw3TexB9wK5"},
	{"6556f0b7-b826-4ed3-b2ef-911a1d0665c8", "RfzKJAwgPxnMD24SEVxj3L"},
	{"269ca19c-5937-4512-83f8-ab9ab82cdaab", "nUkCVCrYTu7LBY8eDPCbs8"},
	{"58f3c051-a71a-4608-bf4e-543694f65f6b", "kKorfq5sQENdZGWDPuC8qH"},
	{"13f0a717-883d-4a69-b7da-eeb91f5d5ca5", "Dke3AiNbfvWUjq6W5bfEZ5"},
	{"7d2ec191-96ef-4b3d-834e-94cf60b1467f", "6mVb7WwQ8tcas48qYJoYHQ"},
	{"831805c0-7651-4a0d-804e-d12924bb9a0c", "pZ6gS9wWUXdHr64W7jrVLR"},
	{"56a03672-74bc-418c-9964-76548fd5b209", "9qiVxV4fEdYvhaW5PBSXRH"},
	{"8c634f54-f7bf-4f39-bf84-394f06970eda", "YpGmB2LeJ6ms8jntjgTkyS"},
	{"a5d35565-a3e0-407f-ab1c-1b8a605ba172", "9bFhyhbykoBC3mdQcu4jWX"},
	{"48443c48-0632-4b72-ba41-f4bdd0cd8ef5", "8ZmSNTgRwuf3qD5ZuLptrE"},
	{"1e86f5b0-b479-4dd2-9d3c-0c0b9a63db30", "WbphNvGKPP8bQE9y65jbS7"},
	{"7178fac6-27b6-473d-801a-63be24d75fb2", "bTdfo4UL4NWDXoMttdcnCN"},
	{"8e03626b-9c23-4032-9ee2-ae6ca661b7d1", "zvvZNCeKsWEQwYrP4dyEHT"},
	{"51f7acb3-2a82-4cc6-8601-523ed6d5909b", "EdQ8N465wYULNARFAfbHbG"},
	{"676fcb0c-730c-4875-9c0b-62b0bc2f4b9e", "J7DMkkwPsVJnt6frUnD2RL"},
	{"2ab244e0-4e34-494f-8c1a-d91051a4c3bf", "E4VoERm26o5Rkfsp6NK2c9"},
	{"4c8c949d-5bef-4d56-82f9-2041859b1dc6", "sEsXtdJXEnq9cLJuCrSLdF"},
	{"ca887b94-e90a-4143-b498-b0f517677e5a", "bLBXoAmYhoJvTHv2FVUz3e"},
	{"eb9ccaea-9569-4688-a40a-96fe3dd3d6eb", "89hHCpRTfBuhi7hxJDWUvj"},
	{"00000000-0000-0000-0000-000000000000", "2222222222222"},
}

func TestGeneration(t *testing.T) {
	tests := []string{
		"",
		"http://www.example.com/",
		"HTTP://www.example.com/",
		"example.com/",
	}

	for _, test := range tests {
		u := NewWithNamespace(test)
		if len(u) < 20 || len(u) > 24 {
			t.Errorf("expected %q to be in range [20, 24], got %d", u, len(u))
		}
	}
}

func TestEncoding(t *testing.T) {
	for _, test := range testVector {
		u, err := uuid.FromString(test.uuid)
		if err != nil {
			t.Error(err)
		}

		suuid := DefaultEncoder.Encode(u)

		if suuid != test.shortuuid {
			t.Errorf("expected %q, got %q", test.shortuuid, suuid)
		}
	}
}

func TestDecoding(t *testing.T) {
	for _, test := range testVector {
		u1, err := uuid.FromString(test.uuid)
		if err != nil {
			t.Error(err)
		}

		u2, err := DefaultEncoder.Decode(test.shortuuid)
		if err != nil {
			t.Error(err)
		}

		if u1 != u2 {
			t.Errorf("expected %q, got %q", u1, u2)
		}
	}
}

func TestNewWithAlphabet(t *testing.T) {
	abc := DefaultAlphabet[:len(DefaultAlphabet)-1] + "="
	enc := base57{newAlphabet(abc)}
	u1, _ := uuid.FromString("e9ae9ba7-4fb1-4a6d-bbca-5315ed438371")
	u2 := enc.Encode(u1)
	if u2 != "u=BFWRLr5dXbeWf==iasZi" {
		t.Errorf("expected uuid to be %q, got %q", "u=BFWRLr5dXbeWf==iasZi", u2)
	}
}

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkEncoding(b *testing.B) {
	u := uuid.Must(uuid.NewV4())
	for i := 0; i < b.N; i++ {
		DefaultEncoder.Encode(u)
	}
}

func BenchmarkDecoding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DefaultEncoder.Decode("c3eeb3e6-e577-4de2-b5bb-08371196b453")
	}
}
