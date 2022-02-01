package shortuuid

import (
	"testing"

	"github.com/google/uuid"
)

var testVector = []struct {
	uuid      string
	shortuuid string
}{
	{"00000000-0000-0000-0000-000000000000", "2222222222222222222222"},
	{"00000000-0000-0000-8000-000000000000", "22222222222TVBRkNB8C7z"},
	{"00000000-0000-0001-0000-000000000000", "22222222222txLqViLENDy"},
	{"0026636a-e9b3-4a88-9c66-bf49d8cad81f", "23XgxsETJcF7vB5N3FZRcB"},
	{"00302322-9d08-4394-a1ac-b25bbab774d4", "23uhkLRaaS2Kf2mH6aw3vz"},
	{"0170eb09-dd9d-4682-a58f-b49eb481a40d", "2Gd4ScVghRnBS8sFucX6EY"},
	{"02f3da16-b995-445a-877b-2319f0c06b41", "2XwmfbG9HcicRRB9t5CQa5"},
	{"07dcc180-1323-4985-8359-296261dd53c8", "3QjwhwGAAWbJruhPuZD7wQ"},
	{"08c79281-f45e-4463-b26d-1dc84ece1a05", "3a4BuxdZV5sXHXKmBFLbgV"},
	{"0a4399e6-7fad-452b-acc8-1d6b33415019", "3q7KRHpCPccxxRbNCQFVzN"},
	{"12a12ca0-0d6e-4b84-bd6f-9715d467ae76", "5Kw9BxeT3wnKQypL6ACLCi"},
	{"13f0a717-883d-4a69-b7da-eeb91f5d5ca5", "5ZEfb5W6qjUWvfbNiA3ekD"},
	{"185c8229-dbd0-4fc1-a485-04a0532a71d0", "6M5UgQmVxBe3ckaVTUUiuh"},
	{"1e86f5b0-b479-4dd2-9d3c-0c0b9a63db30", "7Sbj56y9EQb8PPKGvNhpbW"},
	{"1fc9211b-c1b8-464e-ae01-1c59d1972f47", "7fNDVkBfqRWtkJocjEaBNE"},
	{"21e24536-3b3e-4166-85c4-c29ea151fbbe", "83e7tBeJWmgit6GXhpqZNL"},
	{"22eac9a7-b0b7-4c8b-9e88-4a2c4a48eef4", "8E8Qu3GcbJQivscyhKNMKP"},
	{"241d6b4e-81a7-439a-896d-b645ec223eb5", "8SGoRZ8p6ngvj7FqCdJTuf"},
	{"269ca19c-5937-4512-83f8-ab9ab82cdaab", "8sbCPDe8YBL7uTYrCVCkUn"},
	{"2ab244e0-4e34-494f-8c1a-d91051a4c3bf", "9c2KN6psfkR5o62mREoV4E"},
	{"2c301095-3f99-4fc8-88e3-a2be83fd39fc", "9s9SGreyXZoi9pVptkdJAY"},
	{"2f56ebd5-373f-478f-afab-1142b74b75bc", "AS7NiY8gezxKjbcYM4d7iH"},
	{"35165d02-262d-4e30-8344-c70ee9cc63ee", "BTPywsQCF52Bbxm8ca4BjD"},
	{"36a06807-8c07-4e65-8b99-855b3628c155", "BizmDiAy2soYEpWRwZzuXL"},
	{"37eea0c0-ec18-42e8-a152-c98fee7346fb", "BxFTrJCpqKYJZmesa5mZgQ"},
	{"39439367-6f10-4de8-883f-7ace28c3bc33", "CCkMExUkSjDpBmttEQYbSK"},
	{"3b6aa689-ab23-4ded-932d-090e029a7a99", "Caai3cSxUeNerdbzyHpjP4"},
	{"3bf87639-ba74-463e-9455-bf7bcb8fd898", "CgCvZPhTgVY3GTHFRNYUp4"},
	{"3c9886b3-cabc-4854-8fe3-d394c0901830", "CnYNSJhx3SbpWbUjnWAPLE"},
	{"40e86d7b-77dd-4f55-aec8-29262e028c38", "DZH4Y7T35du34RmbKx8dA3"},
	{"422a2339-0950-4a68-aec1-8c4056feef2a", "Dn2VoLxuksZXw9qeLUqxBS"},
	{"431933fb-14dd-411c-8fc6-5ca57fa7979d", "DwVKuDAHLQmC6cdoExWLja"},
	{"464a919f-54e9-435e-a029-a3b942b07134", "EWrz2L3LzLos9tQRXuwvHT"},
	{"4788fd99-a1f3-468a-bc67-ada5f80016f4", "EjV25GUWNHpwTruYfoUqvQ"},
	{"48443c48-0632-4b72-ba41-f4bdd0cd8ef5", "ErtpLuZ5Dq3fuwRgTNSmZ8"},
	{"4860e725-08f5-4ba1-87e6-1db791f183d8", "Et3Z8brWsainxRVcDk6aK8"},
	{"48610926-01fb-4ef8-ae6b-ba9b45569525", "Et3rE6R9NndYve6EHELcY4"},
	{"4aeff9a2-4ad2-4c85-abd1-95e798b2132a", "FLykSMTR7pGEL3WYSv3uPJ"},
	{"4ef326fe-7886-4054-b427-4cd9f48200bf", "G4gCMsrP2yvqccZiZYFgbi"},
	{"502f56e0-f351-4f42-b395-4dde8b8c19b6", "GHDAmjSyHRZ6wZ9WsKAt8B"},
	{"503c91b0-6a52-4b26-be71-b96c1d4d908b", "GHj3aCtTKuuGgmUoe5q2qf"},
	{"51f7acb3-2a82-4cc6-8601-523ed6d5909b", "GbHbfAFRANLUYw564N8QdE"},
	{"52600a3d-a12a-4041-aa9c-afd5e5923a70", "GfRGXRYgzqMgxXdrgmWhHZ"},
	{"52fa1417-02e9-4656-96ec-5cf44fdcd3c2", "GmX6obVNjxbzgpNmdU2FyR"},
	{"56a03672-74bc-418c-9964-76548fd5b209", "HRXSBP5WahvYdEf4VxViq9"},
	{"582dd24f-beda-42cf-b751-7c7e9827b7d1", "HhHGPYkG7ZSiQZhBfJuhjc"},
	{"58f3c051-a71a-4608-bf4e-543694f65f6b", "Hq8CuPDWGZdNEQs5qfroKk"},
	{"61714a37-a894-4ec5-a569-7e925b46e368", "KMECx4sy6ZUKYfM2SHQ6on"},
	{"63815548-e46f-441a-8bdf-3153db75e56b", "Ki9ZNo5CDmk2vTRW7nrmnQ"},
	{"63ac6e18-e0c9-477d-b08f-039b9538e577", "KjqsNx43BesCQBZgscd3hQ"},
	{"64d1355f-d052-4bd9-83f4-39b93fb1c01f", "KwSysDpxcBU9FNhGkn2dCf"},
	{"6556f0b7-b826-4ed3-b2ef-911a1d0665c8", "L3jxVES42DMnxPgwAJKzfR"},
	{"65656f1e-ae8d-48a3-825d-4499144eefab", "L4KgyDNjm4Y2mX3gr7DSGb"},
	{"6611ce43-fe22-47ef-879e-de3fbbd4d9f6", "LB9uvEKNdUUDxCEyxgWsHH"},
	{"676fcb0c-730c-4875-9c0b-62b0bc2f4b9e", "LR2DnUrf6tnJVsPwkkMD7J"},
	{"67cdb44b-f0fb-47c1-9060-025216c14067", "LUjH6MQibnaB9cuXp6n4Dn"},
	{"68252941-638b-4e0d-807a-3e8cd04f2493", "LYCkenaotvTFnuLvXmfPRi"},
	{"69807fd5-8bcf-4d4c-a534-ff298d4a5a0a", "Lmx5Y2Uuss9CyhDstUa9ef"},
	{"6dfacd19-5cfd-4b95-92be-7e9fb0fe343e", "MaNVyTNhnAXWL64gXxrw6C"},
	{"709a8e5a-6320-4649-b4c2-8aceb28502ae", "N3yNX6ojCSX9jHNDTSbLke"},
	{"7178fac6-27b6-473d-801a-63be24d75fb2", "NCncdttMoXDWN4LU4ofdTb"},
	{"7227d320-fd77-4430-8e9c-ac3f93e07c51", "NKiRwULs6zretmpJypPnrc"},
	{"7273ea84-04a1-4d32-92ff-b0afad3eedb6", "NNjFfBxxypBKeZqSgvBBF6"},
	{"748374a4-43c4-42a4-92be-f70daffe5d58", "NjdUFi8oPRTzeNT2J3976h"},
	{"75a8dcb5-9e11-4fbb-895b-60bb6959b6e7", "NwFzavAPm2qifgZgWGqvt9"},
	{"77eee831-5e48-400b-b1e1-31bc8d28876e", "PMKJYoQxvihyx7KNSZWGnh"},
	{"7b80853c-a9f3-4acd-b735-17d29b3b9122", "PyWJgkLnvQeVMVo6TQmpGg"},
	{"7c24ef71-f4eb-4000-afb6-710c7f9e0942", "Q72ZX4v7KuwxJax3gVMXt9"},
	{"7c37b0a5-9d63-4d91-a79e-85a480fd2348", "Q7kuSGAQ8Sb9PEwomu3sog"},
	{"7d2ec191-96ef-4b3d-834e-94cf60b1467f", "QHYoJYq84sact8QwW7bVm6"},
	{"81fa1f63-389e-44ae-906d-4c1f4a715a61", "R9BHaDtdyYBniPRHQJvVQm"},
	{"84973303-f7ff-47a0-8b47-c54767700622", "Rbg7Q98TJXRPn4MSJTnczk"},
	{"8a95bedb-0592-4329-87be-8affdf5f2e97", "SfUDm49oppsBcpVWb394kc"},
	{"8c634f54-f7bf-4f39-bf84-394f06970eda", "SykTgjtnj8sm6JeL2BmGpY"},
	{"8ec76efb-b891-4ee4-868d-997597e390e1", "TQzg6cKDuuHKSknu3qaodm"},
	{"9206e127-8559-4cc4-877e-5cf342aecd25", "Tyw9TbyXF3eboAw8NjppL7"},
	{"92faac68-0387-49ae-a219-45e7c13fd6a3", "UAbeAtPqE2tbvAPWotK8wR"},
	{"97d87c66-f983-48b7-84c4-1d521e790ecf", "V2wmd4zFcini9qywg8LFYj"},
	{"9cbfe3da-dc50-4d07-9d3a-3d161e5a82b1", "VtgZctJUu6M3Vtcyn9vRhd"},
	{"9da4fb93-9260-46a7-b6df-f04218560fe0", "W4ks8eJyfojFeKqTELVBxj"},
	{"9dc996f5-2dfb-4e60-96a2-f11f59eb71de", "W6DXoPVDRwps6zCRDdRt8B"},
	{"9e0ef764-f98d-4c66-abd8-5d5e7392bacb", "W8xCHJW9uuvGZMNfsJXnsT"},
	{"9e5366f3-6bb8-42db-b999-1dc3612be8e4", "WBfie497U49tydVaVqjdWn"},
	{"a23360e0-1caf-4275-b957-0bcd7aa109d8", "WrxfrCUdiVEVag4Mq58VRm"},
	{"a2952f3b-fba2-4926-b19d-26d550093b93", "WvqXULoEyrV6KXWxnCZaEV"},
	{"a5d35565-a3e0-407f-ab1c-1b8a605ba172", "XWj4ucQdm3CBokybhyhFb9"},
	{"a5fb9251-9373-486b-b9bd-ecf1c7cc6d9c", "XYKuwCFTTf6msjALCKWqQD"},
	{"a61daaac-ea00-4bba-89d9-22b97763a43a", "XZfuL6LuHaDAimv3pqFfCa"},
	{"ac6250d9-6e12-4686-98b2-3e5e8ab11ecb", "YgFKdNGJ4zLqVdQGsrXVZZ"},
	{"b3150123-4e7c-4080-a029-c7437c6513df", "ZsBD7sfh8ey2kazad5mqYJ"},
	{"b6f1e412-da4f-43bf-8a9f-ae4658f5b9ce", "aZMBdoduJguHtHspgUhoE7"},
	{"ba40a392-f1f6-4486-8c27-00b78db13ddc", "b9tCXujbMqGfph4TreZHHi"},
	{"c1544cac-c6ee-40a7-98c8-900a7282ae22", "cQf4Gk23KAeWHk2asCxn34"},
	{"c3d5e2c5-1bc8-4a86-96fb-bf0877c9fc03", "cr5noXJRXDT9bE3tmwNqNU"},
	{"c3eeb3e6-e577-4de2-b5bb-08371196b453", "cs4pz4NyK4A5M4MZGvHKSn"},
	{"c5232425-c634-4653-a367-0eb07235ccfc", "d6HK6dG7Wb2oDAP8FtjX6U"},
	{"c6876847-fda6-4ab5-a2c2-9d8fe9c543bc", "dLPn6TVX7qkyi37ckrf5jB"},
	{"c7c541b0-d024-432e-8b84-dcc011bfb596", "dYyWT3AKg2gpzcNYsKy4ga"},
	{"ca887b94-e90a-4143-b498-b0f517677e5a", "e3zUVF2vHTvJohYmAoXBLb"},
	{"cda34e96-188b-41b2-acac-05cd288c32e1", "ebVFH2PRLZGBpMyu76wkeC"},
	{"cf5c7de8-82d2-475f-a078-825b76d0a287", "etxUA83NZc8uSh8BpnDWzP"},
	{"d22b3035-2213-49a8-8399-a547303f29a3", "fQSLY8tdfjzxGhK68a94K9"},
	{"d2af90cd-1b97-4bec-810b-562406b31470", "fVgFoxGYiHXk5HFoJZqfTa"},
	{"dbc3f7d2-5f4b-4fb9-b473-2ca74cf7e2ef", "h7kuXcwCaAVwKSnUe6iGRW"},
	{"ded85fd3-afe7-45a1-a4d8-90c0a5bc8901", "hezCGCpHvPR7GzWQB9CkDM"},
	{"e63d7da2-7f8d-4954-9c2d-d46b324d10f2", "ixyyBTCRsdXRvsrnwrYb3P"},
	{"e9f17c46-0d3d-45e8-bc36-6c793b59c91c", "jdYcXAByajJkJxLCFk7hXQ"},
	{"eb9ccaea-9569-4688-a40a-96fe3dd3d6eb", "jvUWDJxh7ihuBfTRpCHh98"},
	{"eda9eacb-d698-41af-8790-4f26ecbd3a39", "kJHGx6rK2FJetgWVqErX92"},
	{"ee2d3de2-4ad7-4f32-9071-f8aa4b46dc0b", "kPUojQWnRSWd4pc8CraKrX"},
	{"eff670a0-75f5-4458-afd4-9aa442dcc507", "khbCh2ACgUZ7pPJBxSiVef"},
	{"f1975e43-a054-4243-85b1-14ed27be8330", "kz7eWWYJKPs57Ngfoo54bB"},
	{"f4c03714-aa7e-42e3-91f8-9a7a3bc05f75", "mZA6APQhgAyRM7VXPQAtTV"},
	{"f8221ca8-49e2-44e5-a104-7cec6fb28d48", "nATLgZGodTTho3m3tT286Q"},
	{"f9b2ea09-32c6-426b-bb3b-a269182e99f3", "nSLNv7fA4WVkdDSunaFr7b"},
	{"f9ee01c3-2015-4716-930e-4d5449810833", "nUfojcH2M5j9j3Tk5A8mf7"},
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
		u, err := uuid.Parse(test.uuid)
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
		u1, err := uuid.Parse(test.uuid)
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

func TestDecodingErrors(t *testing.T) {
	var (
		NotPartOfAlphabetError  = "not part of alphabet"
		UUIDLengthOverflowError = "UUID length overflow"
	)
	tests := []struct {
		shortuuid string
		error     string
	}{
		{"yoANdrf88xUXvwbS5GRbMN", UUIDLengthOverflowError},
		{"tWkeanmnjjupCMcjnUsfef", UUIDLengthOverflowError},
		{"1lIO022222222222222222", NotPartOfAlphabetError},
		{"0a6hrgRGNfQ57QMHZdNYAg", NotPartOfAlphabetError},
	}
	for _, test := range tests {
		_, err := DefaultEncoder.Decode(test.shortuuid)
		if err == nil {
			t.Errorf("expected %q error for %q", test.error, test.shortuuid)
		}
	}
}

func TestNewWithAlphabet(t *testing.T) {
	abc := DefaultAlphabet[:len(DefaultAlphabet)-1] + "="
	enc := base57{newAlphabet(abc)}
	u1, _ := uuid.Parse("e9ae9ba7-4fb1-4a6d-bbca-5315ed438371")
	u2 := enc.Encode(u1)
	if u2 != "iZsai==fWebXd5rLRWFB=u" {
		t.Errorf("expected uuid to be %q, got %q", "u=BFWRLr5dXbeWf==iasZi", u2)
	}
}

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkEncoding(b *testing.B) {
	u := uuid.New()
	for i := 0; i < b.N; i++ {
		DefaultEncoder.Encode(u)
	}
}

func BenchmarkDecoding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DefaultEncoder.Decode("c3eeb3e6-e577-4de2-b5bb-08371196b453")
	}
}
