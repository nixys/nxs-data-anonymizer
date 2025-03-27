--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: list_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.list_types (
    integer_type integer,
    numeric_type numeric,
    double_precision_type double precision,
    varchar_type character varying,
    text_type text,
    date_type date,
    time_tz_type time with time zone,
    boolean_type boolean,
    xml_type xml,
    jsonb_type jsonb,
    id bigint NOT NULL
);


ALTER TABLE public.list_types OWNER TO postgres;

--
-- Name: list_types2; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.list_types2 (
    integer_type integer,
    numeric_type numeric,
    double_precision_type double precision,
    varchar_type character varying,
    text_type text,
    date_type date,
    time_tz_type time with time zone,
    boolean_type boolean,
    xml_type xml,
    jsonb_type jsonb,
    id bigint NOT NULL
);


ALTER TABLE public.list_types2 OWNER TO postgres;

--
-- Name: list_types3; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.list_types3 (
    integer_type integer,
    numeric_type numeric,
    double_precision_type double precision,
    varchar_type character varying,
    text_type text,
    date_type date,
    time_tz_type time with time zone,
    boolean_type boolean,
    xml_type xml,
    jsonb_type jsonb,
    id bigint NOT NULL
);


ALTER TABLE public.list_types3 OWNER TO postgres;

--
-- Name: list_types_not_null; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE public.list_types_not_null (
                                   varchar_type character varying NOT NULL,
                                   text_type text NOT NULL,
                                   varchar_type_n character varying,
                                   text_type_n text,
                                   id bigint NOT NULL
);


ALTER TABLE public.list_types_not_null OWNER TO postgres;

--
-- Name: list_types_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.list_types_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.list_types_id_seq OWNER TO postgres;

--
-- Name: list_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.list_types_id_seq OWNED BY public.list_types.id;


--
-- Name: list_types id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.list_types ALTER COLUMN id SET DEFAULT nextval('public.list_types_id_seq'::regclass);


--
-- Data for Name: list_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.list_types (integer_type, numeric_type, double_precision_type, varchar_type, text_type, date_type, time_tz_type, boolean_type, xml_type, jsonb_type, id) FROM stdin;
\N	42.99465	\N	Biba	\N	\N	19:51:50+00	\N	\N	\N	6
\N	-84.46685	\N	Pupa	\N	\N	03:34:36+00	\N	\N	\N	2
\N	72.52040	\N	Lupa	 	\N	15:17:37+00	t	\N	\N	4
\N	99.37111	\N	Boba	 	\N	03:34:36+00	\N	\N	\N	8
\N	-90.90125	\N	Cerebla	\N	\N	22:00:45+00	\N	\N	\N	10
8765542	\N	7.84023409	 	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vel pretium lectus quam id leo in vitae. Dignissim cras tincidunt lobortis feugiat vivamus at augue. Sit amet aliquam id diam maecenas. Ornare lectus sit amet est placerat in egestas erat. Et malesuada fames ac turpis egestas maecenas pharetra convallis. Orci sagittis eu volutpat odio facilisis. Mauris in aliquam sem fringilla ut morbi tincidunt augue interdum. Nisi porta lorem mollis aliquam ut porttitor leo a diam. Amet purus gravida quis blandit turpis cursus in. Risus feugiat in ante metus. Purus viverra accumsan in nisl nisi scelerisque eu ultrices vitae. Faucibus nisl tincidunt eget nullam. Lectus quam id leo in vitae turpis massa. Nisl nunc mi ipsum faucibus vitae aliquet nec ullamcorper sit. Pretium lectus quam id leo in vitae turpis massa. Blandit aliquam etiam erat velit scelerisque in dictum non. Eu nisl nunc mi ipsum faucibus vitae.	1996-06-15	\N	t	<root>\n  <west>721744494.7994413</west>\n  <top>-26197703.515699387</top>\n  <or>event</or>\n  <art>\n    <layers>among</layers>\n    <growth>couple</growth>\n    <off>official</off>\n    <made>spring</made>\n    <poem>1032613367</poem>\n  </art>\n  <strong>198858711</strong>\n</root>	{"truck": 1563128845, "little": -1009907448.0705619, "possible": "gave"}	1
4562892	\N	5.8024375	\N	Massa sapien faucibus et molestie ac. Praesent elementum facilisis leo vel. Turpis egestas pretium aenean pharetra magna. Facilisi cras fermentum odio eu feugiat pretium nibh ipsum consequat. Auctor neque vitae tempus quam pellentesque. Ornare aenean euismod elementum nisi quis eleifend. Purus sit amet luctus venenatis lectus. Tortor consequat id porta nibh venenatis cras sed felis eget. Felis bibendum ut tristique et egestas quis ipsum. Pretium fusce id velit ut tortor pretium viverra. Nam aliquam sem et tortor consequat. Nisl pretium fusce id velit. Sem integer vitae justo eget magna fermentum iaculis. Sit amet nulla facilisi morbi tempus iaculis urna id volutpat. Vitae auctor eu augue ut lectus arcu. Malesuada fames ac turpis egestas. Ac odio tempor orci dapibus ultrices in. Amet massa vitae tortor condimentum lacinia quis vel eros donec. Et malesuada fames ac turpis egestas integer. Sodales ut eu sem integer vitae justo eget.	2031-07-14	\N	f	<root>\n  <melted>\n    <officer>-598731677.6283197</officer>\n    <exist>variety</exist>\n    <sitting>46234734.5525651</sitting>\n    <period>rising</period>\n    <rice>-1556052082.582368</rice>\n  </melted>\n  <got>eleven</got>\n  <onlinetools>1863692621</onlinetools>\n  <this>-139561185.33627748</this>\n  <driver>funny</driver>\n</root>	[{"ago": {"shut": "seven", "found": 469996577.0459976, "climate": "early"}, "scene": "slightly", "medicine": true}, true, true]	3
87689278	\N	-0.31813426	 	Id donec ultrices tincidunt arcu. Id nibh tortor id aliquet lectus. Condimentum mattis pellentesque id nibh tortor id aliquet lectus proin. Cursus vitae congue mauris rhoncus. Eu ultrices vitae auctor eu augue ut lectus arcu bibendum. Sed turpis tincidunt id aliquet. Feugiat scelerisque varius morbi enim nunc faucibus a pellentesque sit. Diam donec adipiscing tristique risus nec feugiat in fermentum posuere. Scelerisque eu ultrices vitae auctor eu augue ut. Volutpat ac tincidunt vitae semper quis.	1976-12-10	\N	t	<root>\n  <mission>rise</mission>\n  <layers>\n    <wire>consist</wire>\n    <dinner>someone</dinner>\n    <electric>white</electric>\n    <heavy>victory</heavy>\n    <ice>-761134241.8632131</ice>\n  </layers>\n  <help>1111641157</help>\n  <comfortable>744054503</comfortable>\n  <highway>1876042205.1417785</highway>\n</root>	{"torn": [false, -1325125579, false], "rising": false, "volume": "grow"}	9
687527896	\N	-1.81914025	 	Facilisi etiam dignissim diam quis enim. Diam ut venenatis tellus in metus vulputate eu. Mattis rhoncus urna neque viverra justo nec ultrices. Sagittis nisl rhoncus mattis rhoncus urna neque viverra. Nec ullamcorper sit amet risus nullam eget felis eget. Fames ac turpis egestas maecenas pharetra convallis posuere. Eget arcu dictum varius duis at consectetur lorem. Porta lorem mollis aliquam ut porttitor leo a diam sollicitudin. Magna fermentum iaculis eu non diam phasellus vestibulum lorem sed. Etiam non quam lacus suspendisse. Parturient montes nascetur ridiculus mus. Ornare suspendisse sed nisi lacus sed viverra tellus in. Interdum velit euismod in pellentesque massa placerat duis ultricies lacus. Urna nec tincidunt praesent semper feugiat nibh sed pulvinar.	2051-05-16	\N	t	<root>\n  <piece>\n    <vertical>398209907.5902777</vertical>\n    <prove>wind</prove>\n    <fair>parallel</fair>\n    <closely>paint</closely>\n    <see>-1386172501</see>\n  </piece>\n  <hide>previous</hide>\n  <among>skill</among>\n  <our>conversation</our>\n  <cell>keep</cell>\n</root>	[[-128477882.56726694, 1569798072.8967233, {"us": 438307927.6373253, "think": -870063635, "product": "stuck"}], [{"name": "broad", "took": "brother", "alphabet": -731303259}, [-1065068182.0998564, "finish", "up"], false], false]	5
8767542	\N	5.08081291	\N	Fames ac turpis egestas maecenas. Volutpat lacus laoreet non curabitur gravida arcu ac tortor. Sit amet commodo nulla facilisi nullam vehicula. Ipsum dolor sit amet consectetur adipiscing elit pellentesque. Maecenas ultricies mi eget mauris pharetra. Sed faucibus turpis in eu mi bibendum. Massa ultricies mi quis hendrerit dolor magna. Non diam phasellus vestibulum lorem sed. Vestibulum mattis ullamcorper velit sed ullamcorper morbi tincidunt ornare massa. Massa ultricies mi quis hendrerit dolor magna eget est lorem. Quam elementum pulvinar etiam non quam lacus suspendisse faucibus.	2058-03-29	\N	f	<root>\n  <invented>concerned</invented>\n  <feet>-1671936329.7007055</feet>\n  <win>slept</win>\n  <until>how</until>\n  <lion>-1838021068</lion>\n</root>	{"basis": "mine", "company": {"tired": false, "prevent": false, "suppose": 735075799}, "worried": {"iron": 378563223.91183805, "nest": false, "raise": {"date": false, "engineer": true, "television": 136840736.65782642}}}	7
\.


--
-- Data for Name: list_types2; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.list_types2 (integer_type, numeric_type, double_precision_type, varchar_type, text_type, date_type, time_tz_type, boolean_type, xml_type, jsonb_type, id) FROM stdin;
8765542	\N	7.84023409	 	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vel pretium lectus quam id leo in vitae. Dignissim cras tincidunt lobortis feugiat vivamus at augue. Sit amet aliquam id diam maecenas. Ornare lectus sit amet est placerat in egestas erat. Et malesuada fames ac turpis egestas maecenas pharetra convallis. Orci sagittis eu volutpat odio facilisis. Mauris in aliquam sem fringilla ut morbi tincidunt augue interdum. Nisi porta lorem mollis aliquam ut porttitor leo a diam. Amet purus gravida quis blandit turpis cursus in. Risus feugiat in ante metus. Purus viverra accumsan in nisl nisi scelerisque eu ultrices vitae. Faucibus nisl tincidunt eget nullam. Lectus quam id leo in vitae turpis massa. Nisl nunc mi ipsum faucibus vitae aliquet nec ullamcorper sit. Pretium lectus quam id leo in vitae turpis massa. Blandit aliquam etiam erat velit scelerisque in dictum non. Eu nisl nunc mi ipsum faucibus vitae.	1996-06-15	\N	t	<root>\n  <west>721744494.7994413</west>\n  <top>-26197703.515699387</top>\n  <or>event</or>\n  <art>\n    <layers>among</layers>\n    <growth>couple</growth>\n    <off>official</off>\n    <made>spring</made>\n    <poem>1032613367</poem>\n  </art>\n  <strong>198858711</strong>\n</root>	{"truck": 1563128845, "little": -1009907448.0705619, "possible": "gave"}	1
4562892	\N	5.8024375	\N	Massa sapien faucibus et molestie ac. Praesent elementum facilisis leo vel. Turpis egestas pretium aenean pharetra magna. Facilisi cras fermentum odio eu feugiat pretium nibh ipsum consequat. Auctor neque vitae tempus quam pellentesque. Ornare aenean euismod elementum nisi quis eleifend. Purus sit amet luctus venenatis lectus. Tortor consequat id porta nibh venenatis cras sed felis eget. Felis bibendum ut tristique et egestas quis ipsum. Pretium fusce id velit ut tortor pretium viverra. Nam aliquam sem et tortor consequat. Nisl pretium fusce id velit. Sem integer vitae justo eget magna fermentum iaculis. Sit amet nulla facilisi morbi tempus iaculis urna id volutpat. Vitae auctor eu augue ut lectus arcu. Malesuada fames ac turpis egestas. Ac odio tempor orci dapibus ultrices in. Amet massa vitae tortor condimentum lacinia quis vel eros donec. Et malesuada fames ac turpis egestas integer. Sodales ut eu sem integer vitae justo eget.	2031-07-14	\N	f	<root>\n  <melted>\n    <officer>-598731677.6283197</officer>\n    <exist>variety</exist>\n    <sitting>46234734.5525651</sitting>\n    <period>rising</period>\n    <rice>-1556052082.582368</rice>\n  </melted>\n  <got>eleven</got>\n  <onlinetools>1863692621</onlinetools>\n  <this>-139561185.33627748</this>\n  <driver>funny</driver>\n</root>	[{"ago": {"shut": "seven", "found": 469996577.0459976, "climate": "early"}, "scene": "slightly", "medicine": true}, true, true]	3
87689278	\N	-0.31813426	 	Id donec ultrices tincidunt arcu. Id nibh tortor id aliquet lectus. Condimentum mattis pellentesque id nibh tortor id aliquet lectus proin. Cursus vitae congue mauris rhoncus. Eu ultrices vitae auctor eu augue ut lectus arcu bibendum. Sed turpis tincidunt id aliquet. Feugiat scelerisque varius morbi enim nunc faucibus a pellentesque sit. Diam donec adipiscing tristique risus nec feugiat in fermentum posuere. Scelerisque eu ultrices vitae auctor eu augue ut. Volutpat ac tincidunt vitae semper quis.	1976-12-10	\N	t	<root>\n  <mission>rise</mission>\n  <layers>\n    <wire>consist</wire>\n    <dinner>someone</dinner>\n    <electric>white</electric>\n    <heavy>victory</heavy>\n    <ice>-761134241.8632131</ice>\n  </layers>\n  <help>1111641157</help>\n  <comfortable>744054503</comfortable>\n  <highway>1876042205.1417785</highway>\n</root>	{"torn": [false, -1325125579, false], "rising": false, "volume": "grow"}	9
687527896	\N	-1.81914025	 	Facilisi etiam dignissim diam quis enim. Diam ut venenatis tellus in metus vulputate eu. Mattis rhoncus urna neque viverra justo nec ultrices. Sagittis nisl rhoncus mattis rhoncus urna neque viverra. Nec ullamcorper sit amet risus nullam eget felis eget. Fames ac turpis egestas maecenas pharetra convallis posuere. Eget arcu dictum varius duis at consectetur lorem. Porta lorem mollis aliquam ut porttitor leo a diam sollicitudin. Magna fermentum iaculis eu non diam phasellus vestibulum lorem sed. Etiam non quam lacus suspendisse. Parturient montes nascetur ridiculus mus. Ornare suspendisse sed nisi lacus sed viverra tellus in. Interdum velit euismod in pellentesque massa placerat duis ultricies lacus. Urna nec tincidunt praesent semper feugiat nibh sed pulvinar.	2051-05-16	\N	t	<root>\n  <piece>\n    <vertical>398209907.5902777</vertical>\n    <prove>wind</prove>\n    <fair>parallel</fair>\n    <closely>paint</closely>\n    <see>-1386172501</see>\n  </piece>\n  <hide>previous</hide>\n  <among>skill</among>\n  <our>conversation</our>\n  <cell>keep</cell>\n</root>	[[-128477882.56726694, 1569798072.8967233, {"us": 438307927.6373253, "think": -870063635, "product": "stuck"}], [{"name": "broad", "took": "brother", "alphabet": -731303259}, [-1065068182.0998564, "finish", "up"], false], false]	5
8767542	\N	5.08081291	\N	Fames ac turpis egestas maecenas. Volutpat lacus laoreet non curabitur gravida arcu ac tortor. Sit amet commodo nulla facilisi nullam vehicula. Ipsum dolor sit amet consectetur adipiscing elit pellentesque. Maecenas ultricies mi eget mauris pharetra. Sed faucibus turpis in eu mi bibendum. Massa ultricies mi quis hendrerit dolor magna. Non diam phasellus vestibulum lorem sed. Vestibulum mattis ullamcorper velit sed ullamcorper morbi tincidunt ornare massa. Massa ultricies mi quis hendrerit dolor magna eget est lorem. Quam elementum pulvinar etiam non quam lacus suspendisse faucibus.	2058-03-29	\N	f	<root>\n  <invented>concerned</invented>\n  <feet>-1671936329.7007055</feet>\n  <win>slept</win>\n  <until>how</until>\n  <lion>-1838021068</lion>\n</root>	{"basis": "mine", "company": {"tired": false, "prevent": false, "suppose": 735075799}, "worried": {"iron": 378563223.91183805, "nest": false, "raise": {"date": false, "engineer": true, "television": 136840736.65782642}}}	7
\.


--
-- Data for Name: list_types3; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.list_types3 (integer_type, numeric_type, double_precision_type, varchar_type, text_type, date_type, time_tz_type, boolean_type, xml_type, jsonb_type, id) FROM stdin;
8765542	\N	7.84023409	 	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vel pretium lectus quam id leo in vitae. Dignissim cras tincidunt lobortis feugiat vivamus at augue. Sit amet aliquam id diam maecenas. Ornare lectus sit amet est placerat in egestas erat. Et malesuada fames ac turpis egestas maecenas pharetra convallis. Orci sagittis eu volutpat odio facilisis. Mauris in aliquam sem fringilla ut morbi tincidunt augue interdum. Nisi porta lorem mollis aliquam ut porttitor leo a diam. Amet purus gravida quis blandit turpis cursus in. Risus feugiat in ante metus. Purus viverra accumsan in nisl nisi scelerisque eu ultrices vitae. Faucibus nisl tincidunt eget nullam. Lectus quam id leo in vitae turpis massa. Nisl nunc mi ipsum faucibus vitae aliquet nec ullamcorper sit. Pretium lectus quam id leo in vitae turpis massa. Blandit aliquam etiam erat velit scelerisque in dictum non. Eu nisl nunc mi ipsum faucibus vitae.	1996-06-15	\N	t	<root>\n  <west>721744494.7994413</west>\n  <top>-26197703.515699387</top>\n  <or>event</or>\n  <art>\n    <layers>among</layers>\n    <growth>couple</growth>\n    <off>official</off>\n    <made>spring</made>\n    <poem>1032613367</poem>\n  </art>\n  <strong>198858711</strong>\n</root>	{"truck": 1563128845, "little": -1009907448.0705619, "possible": "gave"}	1
4562892	\N	5.8024375	\N	Massa sapien faucibus et molestie ac. Praesent elementum facilisis leo vel. Turpis egestas pretium aenean pharetra magna. Facilisi cras fermentum odio eu feugiat pretium nibh ipsum consequat. Auctor neque vitae tempus quam pellentesque. Ornare aenean euismod elementum nisi quis eleifend. Purus sit amet luctus venenatis lectus. Tortor consequat id porta nibh venenatis cras sed felis eget. Felis bibendum ut tristique et egestas quis ipsum. Pretium fusce id velit ut tortor pretium viverra. Nam aliquam sem et tortor consequat. Nisl pretium fusce id velit. Sem integer vitae justo eget magna fermentum iaculis. Sit amet nulla facilisi morbi tempus iaculis urna id volutpat. Vitae auctor eu augue ut lectus arcu. Malesuada fames ac turpis egestas. Ac odio tempor orci dapibus ultrices in. Amet massa vitae tortor condimentum lacinia quis vel eros donec. Et malesuada fames ac turpis egestas integer. Sodales ut eu sem integer vitae justo eget.	2031-07-14	\N	f	<root>\n  <melted>\n    <officer>-598731677.6283197</officer>\n    <exist>variety</exist>\n    <sitting>46234734.5525651</sitting>\n    <period>rising</period>\n    <rice>-1556052082.582368</rice>\n  </melted>\n  <got>eleven</got>\n  <onlinetools>1863692621</onlinetools>\n  <this>-139561185.33627748</this>\n  <driver>funny</driver>\n</root>	[{"ago": {"shut": "seven", "found": 469996577.0459976, "climate": "early"}, "scene": "slightly", "medicine": true}, true, true]	3
87689278	\N	-0.31813426	 	Id donec ultrices tincidunt arcu. Id nibh tortor id aliquet lectus. Condimentum mattis pellentesque id nibh tortor id aliquet lectus proin. Cursus vitae congue mauris rhoncus. Eu ultrices vitae auctor eu augue ut lectus arcu bibendum. Sed turpis tincidunt id aliquet. Feugiat scelerisque varius morbi enim nunc faucibus a pellentesque sit. Diam donec adipiscing tristique risus nec feugiat in fermentum posuere. Scelerisque eu ultrices vitae auctor eu augue ut. Volutpat ac tincidunt vitae semper quis.	1976-12-10	\N	t	<root>\n  <mission>rise</mission>\n  <layers>\n    <wire>consist</wire>\n    <dinner>someone</dinner>\n    <electric>white</electric>\n    <heavy>victory</heavy>\n    <ice>-761134241.8632131</ice>\n  </layers>\n  <help>1111641157</help>\n  <comfortable>744054503</comfortable>\n  <highway>1876042205.1417785</highway>\n</root>	{"torn": [false, -1325125579, false], "rising": false, "volume": "grow"}	9
687527896	\N	-1.81914025	 	Facilisi etiam dignissim diam quis enim. Diam ut venenatis tellus in metus vulputate eu. Mattis rhoncus urna neque viverra justo nec ultrices. Sagittis nisl rhoncus mattis rhoncus urna neque viverra. Nec ullamcorper sit amet risus nullam eget felis eget. Fames ac turpis egestas maecenas pharetra convallis posuere. Eget arcu dictum varius duis at consectetur lorem. Porta lorem mollis aliquam ut porttitor leo a diam sollicitudin. Magna fermentum iaculis eu non diam phasellus vestibulum lorem sed. Etiam non quam lacus suspendisse. Parturient montes nascetur ridiculus mus. Ornare suspendisse sed nisi lacus sed viverra tellus in. Interdum velit euismod in pellentesque massa placerat duis ultricies lacus. Urna nec tincidunt praesent semper feugiat nibh sed pulvinar.	2051-05-16	\N	t	<root>\n  <piece>\n    <vertical>398209907.5902777</vertical>\n    <prove>wind</prove>\n    <fair>parallel</fair>\n    <closely>paint</closely>\n    <see>-1386172501</see>\n  </piece>\n  <hide>previous</hide>\n  <among>skill</among>\n  <our>conversation</our>\n  <cell>keep</cell>\n</root>	[[-128477882.56726694, 1569798072.8967233, {"us": 438307927.6373253, "think": -870063635, "product": "stuck"}], [{"name": "broad", "took": "brother", "alphabet": -731303259}, [-1065068182.0998564, "finish", "up"], false], false]	5
8767542	\N	5.08081291	\N	Fames ac turpis egestas maecenas. Volutpat lacus laoreet non curabitur gravida arcu ac tortor. Sit amet commodo nulla facilisi nullam vehicula. Ipsum dolor sit amet consectetur adipiscing elit pellentesque. Maecenas ultricies mi eget mauris pharetra. Sed faucibus turpis in eu mi bibendum. Massa ultricies mi quis hendrerit dolor magna. Non diam phasellus vestibulum lorem sed. Vestibulum mattis ullamcorper velit sed ullamcorper morbi tincidunt ornare massa. Massa ultricies mi quis hendrerit dolor magna eget est lorem. Quam elementum pulvinar etiam non quam lacus suspendisse faucibus.	2058-03-29	\N	f	<root>\n  <invented>concerned</invented>\n  <feet>-1671936329.7007055</feet>\n  <win>slept</win>\n  <until>how</until>\n  <lion>-1838021068</lion>\n</root>	{"basis": "mine", "company": {"tired": false, "prevent": false, "suppose": 735075799}, "worried": {"iron": 378563223.91183805, "nest": false, "raise": {"date": false, "engineer": true, "television": 136840736.65782642}}}	7
\.


--
-- Data for Name: list_types_not_null; Type: TABLE DATA; Schema: public; Owner: postgres
--
COPY public.list_types_not_null (varchar_type, text_type, varchar_type_n, text_type_n, id) FROM stdin;
Non empty varchar	Non empty text	Non empty varchar	Non empty text	0
	Non empty text	Non empty varchar	Non empty text	1
Non empty varchar		Non empty varchar	Non empty text	2
Non empty varchar	Non empty text	\N	\N	3
	Non empty text	\N	\N	4
		\N	\N	5
\.

--
-- Name: list_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.list_types_id_seq', 10, true);


--
-- Name: list_types list_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.list_types
    ADD CONSTRAINT list_types_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

